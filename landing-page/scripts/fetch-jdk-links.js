#!/usr/bin/env node
/**
 * Fetch JDK download links from various distributors
 * and update data/jdk.json
 */

const fs = require('fs').promises;
const path = require('path');

const PLATFORMS = ['windows-x64', 'linux-x64', 'linux-arm64', 'macos-x64', 'macos-arm64'];
const JDK_VERSIONS = [8, 11, 17, 21];

// Platform mapping for different APIs
const PLATFORM_MAP = {
  'windows-x64': { adoptium: { os: 'windows', arch: 'x64' }, jenv: 'windows-x86_64' },
  'linux-x64': { adoptium: { os: 'linux', arch: 'x64' }, jenv: 'linux-x86_64' },
  'linux-arm64': { adoptium: { os: 'linux', arch: 'aarch64' }, jenv: 'linux-aarch_64' },
  'macos-x64': { adoptium: { os: 'mac', arch: 'x64' }, jenv: 'osx-x86_64' },
  'macos-arm64': { adoptium: { os: 'mac', arch: 'aarch64' }, jenv: 'osx-aarch_64' }
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
 * Fetch JDK from Adoptium (Eclipse Temurin)
 */
async function fetchAdoptiumJDK(version, platform) {
  const { os, arch } = PLATFORM_MAP[platform].adoptium;
  const url = `https://api.adoptium.net/v3/assets/latest/${version}/hotspot?architecture=${arch}&image_type=jdk&os=${os}&vendor=eclipse`;

  try {
    const data = await fetchWithRetry(url);
    if (!data || data.length === 0) {
      console.warn(`No Adoptium JDK found for ${platform} version ${version}`);
      return null;
    }

    const binary = data[0].binary;
    return {
      url: binary.package.link,
      size: formatBytes(binary.package.size),
      sha256: binary.package.checksum
    };
  } catch (error) {
    console.error(`Failed to fetch Adoptium JDK ${version} for ${platform}:`, error.message);
    return null;
  }
}

/**
 * Fetch all Adoptium versions
 */
async function fetchAdoptiumDistribution() {
  console.log('Fetching Adoptium (Temurin) JDK links...');
  const versions = {};

  for (const version of JDK_VERSIONS) {
    versions[version] = {};
    console.log(`  Fetching JDK ${version}...`);

    for (const platform of PLATFORMS) {
      const jdk = await fetchAdoptiumJDK(version, platform);
      if (jdk) {
        versions[version][platform] = jdk;
      }
    }
  }

  return {
    name: 'Eclipse Temurin',
    description: 'Most popular open-source JDK distribution',
    recommended: true,
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
    console.log('Starting JDK links fetch...\n');

    // Fetch JEnv releases
    const jenv = await fetchJenvReleases();
    console.log(`✓ JEnv version ${jenv.version} fetched\n`);

    // Fetch Adoptium JDK
    const temurin = await fetchAdoptiumDistribution();
    console.log('✓ Adoptium JDK fetched\n');

    // Build final JSON
    const output = {
      lastUpdated: new Date().toISOString(),
      jenv,
      jdk: {
        versions: JDK_VERSIONS,
        recommended: 11,
        distributions: {
          temurin
        }
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
    console.log(`  Distributions: 1 (Temurin)`);

  } catch (error) {
    console.error('\n❌ Error:', error.message);
    process.exit(1);
  }
}

main();
