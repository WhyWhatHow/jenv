/**
 * Internationalization (i18n) support
 */

const translations = {
  en: {
    // Meta
    title: 'JEnv - Java Version Manager',
    description: 'Fast and easy Java version management tool',

    // Hero
    detected: 'Detected',
    tagline: 'Install JDK in 3 steps, no manual PATH configuration needed',

    // Quick Start
    quickStart: '🚀 Quick Start',
    downloadJenv: 'Download JEnv',
    downloadJdk: 'Download JDK',
    install: 'Install & Use',
    version: 'Version',
    downloading: 'Loading...',
    download: 'Download',
    notAvailable: 'Not available for your platform',

    // JDK Selection
    selectDistribution: 'Select distribution',
    selectVersion: 'Version',
    jdkHint: '💡 Not sure? JDK 11 is perfect for most beginners',

    // Features
    whyJenv: 'Why JEnv?',
    feature1Title: 'Fast Switching',
    feature1Desc: 'Symlink-based, switch JDK in 300ms',
    feature2Title: 'Windows Optimized',
    feature2Desc: 'Optimized for Windows 10/11, auto privilege handling',
    feature3Title: 'Modern Experience',
    feature3Desc: 'Colorful CLI, light/dark theme support',

    // FAQ
    faqTitle: 'FAQ',

    // Footer
    footerText: 'JEnv - Make Java version management easy',
    license: 'License',
    lastUpdated: 'Last updated',

    // Actions
    copy: 'Copy',
    copied: 'Copied to clipboard',
    copyFailed: 'Copy failed, please copy manually',

    // Errors
    loading: 'Loading...',
    errorLoading: 'Failed to load data. Please refresh the page.',
    initError: 'Initialization failed. Please refresh the page.',

    // Mobile
    mobileTitle: 'Please use desktop browser',
    mobileMessage: 'For the best experience, please open this page on your desktop computer.'
  },

  zh: {
    // Meta
    title: 'JEnv - Java 版本管理工具',
    description: '快速简单的 Java 版本管理工具',

    // Hero
    detected: '检测到',
    tagline: '3 步完成 JDK 安装，再也不用手动配置环境变量',

    // Quick Start
    quickStart: '🚀 快速开始',
    downloadJenv: '下载 JEnv',
    downloadJdk: '下载 JDK',
    install: '安装并使用',
    version: '版本',
    downloading: '加载中...',
    download: '下载',
    notAvailable: '您的平台暂不支持',

    // JDK Selection
    selectDistribution: '选择发行版',
    selectVersion: '版本',
    jdkHint: '💡 不确定选哪个？JDK 11 适合大多数初学者',

    // Features
    whyJenv: '为什么选择 JEnv?',
    feature1Title: '秒级切换',
    feature1Desc: '基于符号链接，300ms 完成 JDK 切换',
    feature2Title: 'Windows 优化',
    feature2Desc: '专为 Windows 10/11 优化，自动处理权限',
    feature3Title: '现代体验',
    feature3Desc: '彩色 CLI，支持亮色/暗色主题',

    // FAQ
    faqTitle: '常见问题',

    // Footer
    footerText: 'JEnv - 让 Java 版本管理更简单',
    license: '开源协议',
    lastUpdated: '最后更新',

    // Actions
    copy: '复制',
    copied: '已复制到剪贴板',
    copyFailed: '复制失败，请手动复制',

    // Errors
    loading: '加载中...',
    errorLoading: '加载数据失败，请刷新页面',
    initError: '初始化失败，请刷新页面',

    // Mobile
    mobileTitle: '请使用桌面浏览器访问',
    mobileMessage: '为了获得最佳体验，请在桌面电脑上打开此页面'
  }
};

/**
 * Get translation
 */
function t(key, lang = 'en') {
  return translations[lang]?.[key] || translations.en[key] || key;
}

/**
 * Get browser language
 */
function getLanguage() {
  const browserLang = navigator.language || navigator.userLanguage;
  return browserLang.toLowerCase().includes('zh') ? 'zh' : 'en';
}

// Export for use in other scripts
if (typeof module !== 'undefined' && module.exports) {
  module.exports = { t, getLanguage, translations };
}
