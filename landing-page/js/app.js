/**
 * JEnv Landing Page - Main Application
 */

// Global state
let appData = null;
let currentPlatform = null;
let currentLang = 'en';

/**
 * Detect user platform
 */
function detectPlatform() {
  const ua = navigator.userAgent.toLowerCase();
  const platform = navigator.platform.toLowerCase();

  // Mobile devices
  if (/android|webos|iphone|ipad|ipod|blackberry|iemobile|opera mini/i.test(ua)) {
    return { os: 'mobile', arch: null };
  }

  // Windows
  if (/windows|win32|win64/i.test(ua) || platform.includes('win')) {
    const arch = /win64|x64|amd64|wow64/i.test(ua) ? 'x64' :
                 /arm64|aarch64/i.test(ua) ? 'arm64' : 'x64';
    return { os: 'windows', arch };
  }

  // macOS
  if (/mac os x|macintosh/i.test(ua) || platform.includes('mac')) {
    // Apple Silicon detection
    const isAppleSilicon = /arm/i.test(ua) ||
                          (navigator.maxTouchPoints > 0 && platform.includes('mac'));
    return { os: 'macos', arch: isAppleSilicon ? 'arm64' : 'x64' };
  }

  // Linux
  if (/linux/i.test(ua) && !/android/i.test(ua)) {
    const arch = /aarch64|armv8|arm64/i.test(ua) ? 'arm64' : 'x64';
    return { os: 'linux', arch };
  }

  // Default to Windows x64
  return { os: 'windows', arch: 'x64' };
}

/**
 * Get platform key
 */
function getPlatformKey(platform) {
  return `${platform.os}-${platform.arch}`;
}

/**
 * Get platform display name
 */
function getPlatformName(platform, lang) {
  const names = {
    'en': {
      'windows-x64': 'Windows (x64)',
      'windows-arm64': 'Windows ARM',
      'linux-x64': 'Linux (x64)',
      'linux-arm64': 'Linux ARM',
      'macos-x64': 'macOS (Intel)',
      'macos-arm64': 'macOS (Apple Silicon)'
    },
    'zh': {
      'windows-x64': 'Windows (64位)',
      'windows-arm64': 'Windows ARM',
      'linux-x64': 'Linux (64位)',
      'linux-arm64': 'Linux ARM',
      'macos-x64': 'macOS (Intel)',
      'macos-arm64': 'macOS (Apple Silicon)'
    }
  };
  const key = getPlatformKey(platform);
  return names[lang]?.[key] || names.en[key] || key;
}

/**
 * Load data from JSON
 */
async function loadData() {
  showLoading(true);
  try {
    const response = await fetch('data/jdk.json');
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Failed to load data:', error);
    showError(t('errorLoading', currentLang));
    throw error;
  } finally {
    showLoading(false);
  }
}

/**
 * Render mobile view
 */
function renderMobile() {
  document.body.innerHTML = `
    <div class="mobile-message">
      <div class="mobile-icon">📱</div>
      <h1>${t('mobileTitle', currentLang)}</h1>
      <p>${t('mobileMessage', currentLang)}</p>
      <a href="https://github.com/WhyWhatHow/jenv" class="mobile-btn">
        View on GitHub
      </a>
    </div>
  `;
}

/**
 * Render page
 */
function renderPage() {
  // Update platform info
  document.getElementById('platform-info').textContent =
    getPlatformName(currentPlatform, currentLang);

  // Update last updated timestamp
  if (appData.lastUpdated) {
    const date = new Date(appData.lastUpdated);
    document.getElementById('data-timestamp').textContent = date.toLocaleString();
  }

  // Render JEnv download
  renderJenvDownload();

  // Render JDK download
  renderJdkDownload();

  // Render install commands
  renderInstallCommands();

  // Render FAQ
  renderFAQ();
}

/**
 * Render JEnv download section
 */
function renderJenvDownload() {
  const platformKey = getPlatformKey(currentPlatform);
  const jenvInfo = appData.jenv.platforms[platformKey];

  const versionEl = document.getElementById('jenv-version');
  const downloadBtn = document.getElementById('jenv-download');
  const sizeEl = document.getElementById('jenv-size');

  if (jenvInfo && jenvInfo.url) {
    versionEl.textContent = `v${appData.jenv.version}`;
    downloadBtn.disabled = false;
    downloadBtn.onclick = () => {
      trackDownload('jenv', platformKey);
      window.open(jenvInfo.url, '_blank');
    };
    downloadBtn.querySelector('span:last-child').textContent =
      t('download', currentLang) + ` (${jenvInfo.size})`;
    sizeEl.textContent = '';
  } else {
    versionEl.textContent = appData.jenv.version;
    downloadBtn.disabled = true;
    downloadBtn.querySelector('span:last-child').textContent =
      t('notAvailable', currentLang);
  }
}

/**
 * Render JDK download section
 */
function renderJdkDownload() {
  populateDistributionSelector();
  populateVersionSelector();
  setupJdkDownloadHandlers();
}

/**
 * Populate distribution selector
 */
function populateDistributionSelector() {
  const selector = document.getElementById('dist-selector');
  selector.innerHTML = '';

  const distributions = Object.entries(appData.jdk.distributions);

  distributions.forEach(([id, dist]) => {
    const option = document.createElement('option');
    option.value = id;
    option.textContent = dist.name;
    if (dist.recommended) {
      option.textContent += ' ⭐';
      option.selected = true;
    }
    selector.appendChild(option);
  });
}

/**
 * Populate version selector
 */
function populateVersionSelector() {
  const selector = document.getElementById('jdk-version');
  selector.innerHTML = '';

  const versions = appData.jdk.versions;
  const recommended = Array.isArray(appData.jdk.recommended)
    ? appData.jdk.recommended
    : [appData.jdk.recommended];

  versions.forEach(version => {
    const option = document.createElement('option');
    option.value = version;
    option.textContent = `JDK ${version}`;

    if (recommended.includes(version)) {
      option.textContent += ' ⭐';
      // Select first recommended version
      if (!selector.value && recommended.includes(version)) {
        option.selected = true;
      }
    }

    selector.appendChild(option);
  });
}

/**
 * Setup JDK download handlers
 */
function setupJdkDownloadHandlers() {
  const distSelector = document.getElementById('dist-selector');
  const versionSelector = document.getElementById('jdk-version');
  const downloadBtn = document.getElementById('jdk-download');

  const updateDownloadButton = () => {
    const dist = distSelector.value;
    const version = versionSelector.value;
    const platformKey = getPlatformKey(currentPlatform);

    const jdkInfo = appData.jdk.distributions[dist]?.versions[version]?.[platformKey];

    if (jdkInfo && jdkInfo.url) {
      downloadBtn.disabled = false;
      downloadBtn.onclick = () => {
        trackDownload('jdk', platformKey, dist, version);
        window.open(jdkInfo.url, '_blank');
      };

      const distName = appData.jdk.distributions[dist]?.name || dist;
      let btnText = `${t('download', currentLang)} ${distName} ${version}`;
      if (jdkInfo.size) {
        btnText += ` (${jdkInfo.size})`;
      }
      downloadBtn.querySelector('span:last-child').textContent = btnText;
    } else {
      downloadBtn.disabled = true;
      downloadBtn.querySelector('span:last-child').textContent =
        t('notAvailable', currentLang);
    }
  };

  distSelector.addEventListener('change', updateDownloadButton);
  versionSelector.addEventListener('change', updateDownloadButton);

  // Initial update
  updateDownloadButton();
}

/**
 * Render install commands
 */
function renderInstallCommands() {
  const commands = getInstallCommands(currentPlatform.os, currentLang);
  document.getElementById('install-commands').textContent = commands;
}

/**
 * Get install commands for platform
 */
function getInstallCommands(os, lang) {
  const commands = {
    'windows': {
      'en': `# 1. Extract jenv.zip to any directory
# 2. Run PowerShell/CMD as Administrator

jenv init

# 3. Add downloaded JDK (change to actual path)
jenv add jdk11 "C:\\path\\to\\jdk"

# 4. Switch to JDK 11
jenv use jdk11

# 5. Verify installation
java -version`,
      'zh': `# 1. 解压 jenv.zip 到任意目录
# 2. 以管理员身份运行 PowerShell/CMD

jenv init

# 3. 添加下载的 JDK (修改为实际路径)
jenv add jdk11 "C:\\path\\to\\jdk"

# 4. 切换到 JDK 11
jenv use jdk11

# 5. 验证安装
java -version`
    },
    'linux': {
      'en': `# 1. Extract jenv.zip
tar -xzf jenv-*.zip

# 2. Initialize (may need sudo)
./jenv init

# 3. Add JDK
./jenv add jdk11 /path/to/jdk

# 4. Switch version
./jenv use jdk11

# 5. Reload shell
source ~/.bashrc  # or ~/.zshrc

# 6. Verify
java -version`,
      'zh': `# 1. 解压 jenv.zip
tar -xzf jenv-*.zip

# 2. 初始化 (可能需要 sudo)
./jenv init

# 3. 添加 JDK
./jenv add jdk11 /path/to/jdk

# 4. 切换版本
./jenv use jdk11

# 5. 重新加载 shell
source ~/.bashrc  # 或 ~/.zshrc

# 6. 验证
java -version`
    },
    'macos': {
      'en': `# macOS support coming soon
# Stay tuned...`,
      'zh': `# macOS 支持即将推出
# 敬请期待...`
    }
  };

  return commands[os]?.[lang] || commands.windows.en;
}

/**
 * Render FAQ section
 */
function renderFAQ() {
  const faqData = getFAQItems(currentLang);
  const accordion = document.getElementById('faq-accordion');

  accordion.innerHTML = faqData.map((item, index) => `
    <div class="faq-item ${index === 0 ? 'open' : ''}">
      <button class="faq-question" onclick="toggleFAQ(this)">
        <span>${item.question}</span>
        <span class="arrow">▼</span>
      </button>
      <div class="faq-answer">
        ${item.answer}
      </div>
    </div>
  `).join('');
}

/**
 * Get FAQ items
 */
function getFAQItems(lang) {
  const items = {
    'en': [
      {
        question: 'What is JDK? Do I need it?',
        answer: '<p>JDK (Java Development Kit) is a software development kit for Java. You need it to write and run Java programs.</p>'
      },
      {
        question: 'Which JDK version should I choose?',
        answer: '<ul><li><strong>JDK 8</strong>: Legacy projects</li><li><strong>JDK 11</strong>: Recommended for beginners ⭐</li><li><strong>JDK 17</strong>: Modern features</li><li><strong>JDK 21</strong>: Latest LTS</li></ul>'
      },
      {
        question: 'Why do I need JEnv?',
        answer: '<p>If you need to:</p><ul><li>Keep multiple Java versions</li><li>Quickly switch JDK for different projects</li><li>Avoid manual PATH configuration</li></ul><p>JEnv makes it super easy.</p>'
      },
      {
        question: 'Which OS are supported?',
        answer: '<p>Currently:</p><ul><li>✅ Windows 10/11</li><li>✅ Linux (various distros)</li><li>🚧 macOS (coming soon)</li></ul>'
      }
    ],
    'zh': [
      {
        question: '什么是 JDK? 我必须安装吗?',
        answer: '<p>JDK (Java Development Kit) 是 Java 开发工具包，包含编译器、运行环境等。如果您要学习 Java 编程或运行 Java 程序，就必须安装 JDK。</p>'
      },
      {
        question: '我应该选择哪个 JDK 版本?',
        answer: '<ul><li><strong>JDK 8</strong>: 课程使用 Java 8</li><li><strong>JDK 11</strong>: 推荐新手，稳定 ⭐</li><li><strong>JDK 17</strong>: 需要新特性</li><li><strong>JDK 21</strong>: 最新 LTS</li></ul>'
      },
      {
        question: '为什么需要 JEnv?',
        answer: '<p>如果您需要:</p><ul><li>保留多个 Java 版本</li><li>快速切换不同项目的 JDK</li><li>避免手动修改环境变量</li></ul><p>JEnv 让这些操作变得非常简单。</p>'
      },
      {
        question: '支持哪些操作系统?',
        answer: '<p>目前支持:</p><ul><li>✅ Windows 10/11</li><li>✅ Linux (多发行版)</li><li>🚧 macOS (即将支持)</li></ul>'
      }
    ]
  };

  return items[lang] || items.en;
}

/**
 * Toggle FAQ item
 */
function toggleFAQ(button) {
  const item = button.parentElement;
  const isOpen = item.classList.contains('open');

  // Close all items
  document.querySelectorAll('.faq-item').forEach(i => i.classList.remove('open'));

  // Open clicked item if it was closed
  if (!isOpen) {
    item.classList.add('open');
  }
}

/**
 * Copy code to clipboard
 */
function copyCode() {
  const code = document.getElementById('install-commands').textContent;
  navigator.clipboard.writeText(code).then(() => {
    showToast(t('copied', currentLang), 'success');
  }).catch(() => {
    showToast(t('copyFailed', currentLang), 'error');
  });
}

/**
 * Show toast notification
 */
function showToast(message, type = 'info') {
  const toast = document.getElementById('toast');
  toast.textContent = message;
  toast.className = `toast toast-${type} show`;

  setTimeout(() => {
    toast.classList.remove('show');
  }, 3000);
}

/**
 * Show/hide loading overlay
 */
function showLoading(show) {
  document.getElementById('loading').style.display = show ? 'flex' : 'none';
}

/**
 * Show error message
 */
function showError(message) {
  showToast(message, 'error');
}

/**
 * Track download event
 */
function trackDownload(file, platform, dist = null, version = null) {
  // Simple analytics - can be extended
  console.log('Download:', { file, platform, dist, version, timestamp: Date.now() });

  // Could send to analytics service here
  // Example: gtag('event', 'download', { file, platform });
}

/**
 * Setup event listeners
 */
function setupEventListeners() {
  // Copy code button
  document.getElementById('copy-code').addEventListener('click', copyCode);

  // Language toggle
  document.getElementById('lang-toggle').addEventListener('click', toggleLanguage);
}

/**
 * Toggle language
 */
function toggleLanguage() {
  currentLang = currentLang === 'en' ? 'zh' : 'en';
  document.getElementById('lang-toggle').textContent = currentLang === 'en' ? '中文' : 'EN';
  document.documentElement.lang = currentLang;

  // Re-translate page
  translatePage();

  // Re-render dynamic content
  renderPage();
}

/**
 * Translate page
 */
function translatePage() {
  document.querySelectorAll('[data-i18n]').forEach(el => {
    const key = el.getAttribute('data-i18n');
    el.textContent = t(key, currentLang);
  });

  // Update title
  document.title = t('title', currentLang);
}

/**
 * Initialize application
 */
async function init() {
  try {
    // Detect platform
    currentPlatform = detectPlatform();

    // Check for mobile
    if (currentPlatform.os === 'mobile') {
      renderMobile();
      return;
    }

    // Detect language
    currentLang = getLanguage();
    document.documentElement.lang = currentLang;

    // Load data
    appData = await loadData();

    // Translate page
    translatePage();

    // Render page
    renderPage();

    // Setup event listeners
    setupEventListeners();

  } catch (error) {
    console.error('Initialization error:', error);
    showError(t('initError', currentLang));
  }
}

// Start app when DOM is ready
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', init);
} else {
  init();
}
