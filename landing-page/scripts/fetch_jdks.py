#!/usr/bin/env python
import json
import os
import requests
from datetime import datetime, timezone

PLATFORMS = ['windows-x64', 'linux-x64', 'linux-arm64', 'macos-x64', 'macos-arm64', 'windows-arm64']
DISTRIBUTIONS = [
    {'id': 'temurin', 'name': 'Eclipse Temurin', 'desc': 'Most popular open-source JDK', 'recommended': True},
    {'id': 'zulu', 'name': 'Azul Zulu', 'desc': 'Enterprise-ready OpenJDK', 'recommended': True},
    {'id': 'graalvm', 'name': 'GraalVM', 'desc': 'High-performance JDK', 'recommended': True},
    {'id': 'dragonwell', 'name': 'Alibaba Dragonwell', 'desc': 'Alibaba\'s OpenJDK', 'recommended': True},
    {'id': 'oracle_open_jdk', 'name': 'Oracle OpenJDK', 'desc': 'Official OpenJDK builds', 'recommended': True},
    {'id': 'corretto', 'name': 'Amazon Corretto', 'desc': 'Production-ready OpenJDK'},
    {'id': 'liberica', 'name': 'BellSoft Liberica', 'desc': 'Flexible OpenJDK builds' },
    {'id': 'microsoft', 'name': 'Microsoft Build of OpenJDK', 'desc': 'Microsoft\'s OpenJDK' },
    {'id': 'sapmachine', 'name': 'SapMachine', 'desc': 'SAP\'s OpenJDK' }
]
PLATFORM_MAP = {
    'windows-x64': {'foojay': {'os': 'windows', 'arch': 'x64'}, 'jenv': 'windows-x86_64'},
    'linux-x64': {'foojay': {'os': 'linux', 'arch': 'x64'}, 'jenv': 'linux-x86_64'},
    'linux-arm64': {'foojay': {'os': 'linux', 'arch': 'aarch64'}, 'jenv': 'linux-aarch_64'},
    'macos-x64': {'foojay': {'os': 'macos', 'arch': 'x64'}, 'jenv': 'osx-x86_64'},
    'macos-arm64': {'foojay': {'os': 'macos', 'arch': 'aarch64'}, 'jenv': 'osx-aarch_64'},
    'windows-arm64': {'foojay': {'os': 'windows', 'arch': 'aarch64'}, 'jenv': 'windows-aarch_64'}
}

def fetch_with_retry(url, headers=None, retries=3):
    for i in range(retries):
        try:
            response = requests.get(url, headers=headers)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            print(f"Attempt {i + 1} failed for {url}: {e}")
            if i == retries - 1:
                raise
            import time
            time.sleep(1 * (i + 1))

def fetch_jenv_releases():
    print('Fetching JEnv releases...')
    url = 'https://api.github.com/repos/WhyWhatHow/jenv/releases/latest'
    headers = {'User-Agent': 'jenv-landing-fetcher'}
    github_token = os.getenv('GITHUB_TOKEN')
    if github_token:
        print("Found GITHUB_TOKEN, using it for authentication.")
        headers['Authorization'] = f'Bearer {github_token}'
    data = fetch_with_retry(url, headers=headers)
    version = data['tag_name'].replace('v', '')
    platforms = {}
    for platform in PLATFORMS:
        if platform not in PLATFORM_MAP:
            continue
        platform_key = PLATFORM_MAP[platform]['jenv']
        asset = next((a for a in data['assets'] if platform_key in a['name']), None)
        if asset:
            platforms[platform] = {
                'url': asset['browser_download_url'],
                'size': format_bytes(asset['size']),
                'sha256': ''
            }
    return {'version': version, 'platforms': platforms}

def fetch_maintained_jdk_versions():
    print('Fetching maintained JDK versions...')
    url = 'https://api.foojay.io/disco/v3.0/major_versions?maintained=true'
    data = fetch_with_retry(url)
    return [int(v['major_version']) for v in data['result'] if int(v['major_version']) >= 8]

def fetch_foojay_jdk(distribution, version, platform):
    foojay_os = PLATFORM_MAP[platform]['foojay']['os']
    foojay_arch = PLATFORM_MAP[platform]['foojay']['arch']

    params = {
        'version': version,
        'distribution': distribution,
        'operating_system': foojay_os,
        'architecture': foojay_arch,
        'archive_type': 'zip' if foojay_os == 'windows' else 'tar.gz',
        'package_type': 'jdk',
        'latest': 'available',
        'release_status': 'ga'
    }

    url = f"https://api.foojay.io/disco/v3.0/packages"
    try:
        response = requests.get(url, params=params)
        response.raise_for_status()
        data = response.json()

        if not data or not data.get('result'):
            print(f"No Foojay JDK found for {distribution} {version} on {platform}")
            return None

        pkg = data['result'][0]
        return {
            'url': pkg['links']['pkg_download_redirect'],
            'size': format_bytes(pkg.get('size', 0)),
            'sha256': pkg.get('checksum', ''),
            'javaVersion': pkg['java_version'],
            'distribution': pkg['distribution']
        }
    except requests.exceptions.RequestException as e:
        print(f"Failed to fetch Foojay {distribution} JDK {version} for {platform}: {e}")
        return None

def format_bytes(bytes_):
    if bytes_ < 1024:
        return f"{bytes_} B"
    if bytes_ < 1024 ** 2:
        return f"{(bytes_ / 1024):.1f} KB"
    return f"{(bytes_ / 1024 ** 2):.1f} MB"

def main():
    try:
        print('Starting JDK links fetch from Foojay DiscoAPI...\n')
        jenv = fetch_jenv_releases()
        print(f"✓ JEnv version {jenv['version']} fetched\n")

        jdk_versions = fetch_maintained_jdk_versions()
        print(f"✓ Maintained JDK versions: {jdk_versions}\n")

        distributions_data = {}
        for dist in DISTRIBUTIONS:
            print(f"Fetching {dist['name']}...")
            versions = {}
            for version in jdk_versions:
                print(f"  JDK {version}...")
                platforms_data = {}
                for platform in PLATFORMS:
                    if platform not in PLATFORM_MAP:
                        continue
                    jdk = fetch_foojay_jdk(dist['id'], str(version), platform)
                    if jdk:
                        platforms_data[platform] = jdk
                    import time
                    time.sleep(0.1)
                if platforms_data:
                    versions[version] = platforms_data

            distributions_data[dist['id']] = {
                'name': dist['name'],
                'description': dist['desc'],
                'recommended': dist.get('recommended', False),
                'versions': versions
            }
            print(f"✓ {dist['name']} fetched\n")

        output = {
            'lastUpdated': datetime.now(timezone.utc).isoformat(),
            'jenv': jenv,
            'jdk': {
                'versions': jdk_versions,
                'recommended': [17, 21, 25],
                'distributions': distributions_data
            }
        }

        output_path = os.path.join(os.path.dirname(__file__), '../data/jdk.json')
        with open(output_path, 'w') as f:
            json.dump(output, f, indent=2)

        print('✓ data/jdk.json updated successfully')

    except Exception as e:
        print(f"\n❌ Error: {e}")
        import traceback
        traceback.print_exc()
        exit(1)

if __name__ == '__main__':
    main()
