#!/usr/bin/env node
/**
 * Fetch JDK download links from Foojay DiscoAPI
 * and update data/jdk.json
 */

const fs = require('fs').promises;
const path = require('path');

const PLATFORMS = ['windows-x64', 'linux-x64', 'linux-arm64', 'macos-x64', 'macos-arm64'];
const JDK_VERSIONS = [8, 11, 17, 21, 25]; // LTS + latest versions

// JDK distributions to fetch (from Foojay)
// For testing, limit to 3 main distributions
const DISTRIBUTIONS = [
  { id: 'temurin', name: 'Eclipse Temurin', desc: 'Most popular open-source JDK', recommended: true },
  { id: 'zulu', name: 'Azul Zulu', desc: 'Enterprise-ready OpenJDK' },
  { id: 'corretto', name: 'Amazon Corretto', desc: 'Production-ready OpenJDK' }
  // Uncomment below for full list
  // { id: 'liberica', name: 'BellSoft Liberica', desc: 'Flexible OpenJDK builds' },
  // { id: 'microsoft', name: 'Microsoft Build of OpenJDK', desc: 'Microsoft\'s OpenJDK' },
  // { id: 'oracle_open_jdk', name: 'Oracle OpenJDK', desc: 'Official OpenJDK builds' },
  // { id: 'graalvm_ce17', name: 'GraalVM CE 17', desc: 'High-performance runtime' },
  // { id: 'graalvm_ce21', name: 'GraalVM CE 21', desc: 'High-performance runtime' },
  // { id: 'sapmachine', name: 'SapMachine', desc: 'SAP\'s OpenJDK' },
  // { id: 'dragonwell', name: 'Alibaba Dragonwell', desc: 'Alibaba\'s OpenJDK' }
];

// Platform mapping for different APIs
const PLATFORM_MAP = {
  'windows-x64': { foojay: { os: 'windows', arch: 'x64' }, jenv: 'windows-x86_64' },
  'linux-x64': { foojay: { os: 'linux', arch: 'x64' }, jenv: 'linux-x86_64' },
  'linux-arm64': { foojay: { os: 'linux', arch: 'aarch64' }, jenv: 'linux-aarch_64' },
  'macos-x64': { foojay: { os: 'macos', arch: 'x64' }, jenv: 'osx-x86_64' },
  'macos-arm64': { foojay: { os: 'macos', arch: 'aarch64' }, jenv: 'osx-aarch_64' }
};

/**
 * Fetch with retry logic
 */
async function fetchWithRetry(url, options = {}, retries = 3) {
  for (let i = 0; i < retries; i++) {
    try {
      const response = await fetch(url, options);
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
      return await response.json();
    } catch (error) {
      console.error(`Attempt ${i + 1} failed for ${url}:`, error.message);
      if (i === retries - 1) throw error;
      await new Promise(resolve => setTimeout(resolve, 1000 * (i + 1)));
    }
  }
}

/**
 * Fetch JEnv releases from GitHub
 */
async function fetchJenvReleases() {
  console.log('Fetching JEnv releases...');
  const url = 'https://api.github.com/repos/WhyWhatHow/jenv/releases/latest';
  const data = await fetchWithRetry(url, {
    headers: {
      'User-Agent': 'jenv-landing-fetcher',
      'Accept': 'application/vnd.github.v3+json'
    }
  });

  const version = data.tag_name.replace('v', '');
  const platforms = {};

  for (const platform of PLATFORMS) {
    const platformKey = PLATFORM_MAP[platform].jenv;
    const asset = data.assets.find(a => a.name.includes(platformKey));

    if (asset) {
      platforms[platform] = {
        url: asset.browser_download_url,
        size: formatBytes(asset.size),
        sha256: '' // GitHub doesn't provide SHA256 in API
      };
    }
  }

  return { version, platforms };
}

/**
 * Fetch JDK from Foojay DiscoAPI
 */
async function fetchFoojayJDK(distribution, version, platform) {
  const { os, arch } = PLATFORM_MAP[platform].foojay;

  // Build query parameters
  const params = new URLSearchParams({
    version: version,
    distribution: distribution,
    operating_system: os,
    architecture: arch,
    archive_type: os === 'windows' ? 'zip' : 'tar.gz',
    package_type: 'jdk',
    latest: 'available',
    release_status: 'ga'
  });

  const url = `https://api.foojay.io/disco/v3.0/packages?${params}`;

  try {
    const data = await fetchWithRetry(url);
    if (!data || !data.result || data.result.length === 0) {
      console.warn(`No Foojay JDK found for ${distribution} ${version} on ${platform}`);
      return null;
    }

    // Get the first (latest) package
    const pkg = data.result[0];
    return {
      url: pkg.links?.pkg_download_redirect || pkg.filename,
      size: formatBytes(pkg.size || 0),
      sha256: pkg.checksum || '',
      javaVersion: pkg.java_version,
      distribution: pkg.distribution
    };
  } catch (error) {
    console.error(`Failed to fetch Foojay ${distribution} JDK ${version} for ${platform}:`, error.message);
    return null;
  }
}

/**
 * Fetch all versions for a distribution
 */
async function fetchDistributionData(dist) {
  console.log(`Fetching ${dist.name}...`);
  const versions = {};

  for (const version of JDK_VERSIONS) {
    // Skip incompatible versions
    if (dist.id === 'graalvm_ce17' && version !== 17) continue;
    if (dist.id === 'graalvm_ce21' && version !== 21) continue;

    console.log(`  JDK ${version}...`);
    versions[version] = {};

    for (const platform of PLATFORMS) {
      const jdk = await fetchFoojayJDK(dist.id, version, platform);
      if (jdk) {
        versions[version][platform] = jdk;
      }
      // Small delay to avoid rate limiting
      await new Promise(resolve => setTimeout(resolve, 100));
    }
  }

  return {
    name: dist.name,
    description: dist.desc,
    recommended: dist.recommended || false,
    versions
  };
}

/**
 * Format bytes to human-readable size
 */
function formatBytes(bytes) {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1048576) return (bytes / 1024).toFixed(1) + ' KB';
  return (bytes / 1048576).toFixed(1) + ' MB';
}

/**
 * Main function
 */
async function main() {
  try {
    console.log('Starting JDK links fetch from Foojay DiscoAPI...\n');

    // Fetch JEnv releases
    const jenv = await fetchJenvReleases();
    console.log(`✓ JEnv version ${jenv.version} fetched\n`);

    // Fetch JDK distributions
    const distributions = {};
    for (const dist of DISTRIBUTIONS) {
      const distData = await fetchDistributionData(dist);
      distributions[dist.id] = distData;
      console.log(`✓ ${dist.name} fetched\n`);
    }

    // Build final JSON
    const output = {
      lastUpdated: new Date().toISOString(),
      jenv,
      jdk: {
        versions: JDK_VERSIONS,
        recommended: [11, 17],
        distributions
      }
    };

    // Write to file
    const outputPath = path.join(__dirname, '../data/jdk.json');
    await fs.writeFile(outputPath, JSON.stringify(output, null, 2));
    console.log('✓ data/jdk.json updated successfully');

    // Print summary
    console.log('\nSummary:');
    console.log(`  JEnv version: ${jenv.version}`);
    console.log(`  JEnv platforms: ${Object.keys(jenv.platforms).length}`);
    console.log(`  JDK versions: ${JDK_VERSIONS.length}`);
    console.log(`  Distributions: ${Object.keys(distributions).length}`);

    // Print distribution stats
    console.log('\nDistributions:');
    for (const [id, dist] of Object.entries(distributions)) {
      const versionCount = Object.keys(dist.versions).length;
      console.log(`  - ${dist.name}: ${versionCount} versions`);
    }

  } catch (error) {
    console.error('\n❌ Error:', error.message);
    console.error(error.stack);
    process.exit(1);
  }
}

main();
