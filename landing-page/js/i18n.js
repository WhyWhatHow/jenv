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

    // Typing effect messages
    typingMessages: [
      '300ms to switch JDK versions',
      'Zero manual PATH configuration',
      'Works on Windows, macOS & Linux'
    ],

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
    selectDistributions: 'Select distributions',
    selectVersion: 'Version',
    selectVersions: 'Select versions',
    yourPlatform: 'Your Platform',
    autoDetected: 'Auto-detected',
    jdkHint: '💡 Not sure? JDK 11 is perfect for most beginners',
    distHint: '💡 Select one or more JDK distributions',
    versionHint: '💡 Recommended: JDK 11 or 17',
    selected: 'Selected',
    packages: 'packages',
    downloadSelected: 'Download Selected',

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
    downloadStarted: 'Downloads started...',

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

    // Typing effect messages
    typingMessages: [
      '300ms 切换 JDK 版本',
      '无需手动配置环境变量',
      '支持 Windows, macOS 和 Linux'
    ],

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
    selectDistributions: '选择发行版',
    selectVersion: '版本',
    selectVersions: '选择版本',
    yourPlatform: '您的平台',
    autoDetected: '自动检测',
    jdkHint: '💡 不确定选哪个？JDK 11 适合大多数初学者',
    distHint: '💡 可以选择一个或多个 JDK 发行版',
    versionHint: '💡 推荐：JDK 11 或 17',
    selected: '已选择',
    packages: '个包',
    downloadSelected: '下载所选',

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
    downloadStarted: '开始下载...',

    // Errors
    loading: '加载中...',
    errorLoading: '加载数据失败，请刷新页面',
    initError: '初始化失败，请刷新页面',

    // Mobile
    mobileTitle: '请使用桌面浏览器访问',
    mobileMessage: '为了获得最佳体验，请在桌面电脑上打开此页面'
  },

  jp: {
    // Meta
    title: 'JEnv - Java バージョン管理ツール',
    description: '高速で簡単な Java バージョン管理ツール',

    // Hero
    detected: '検出',
    tagline: '3ステップでJDKをインストール、手動のパス設定は不要',

    // Typing effect messages
    typingMessages: [
      '300ms で JDK バージョンを切り替え',
      '手動の環境変数設定は不要',
      'Windows、macOS、Linux で動作'
    ],

    // Quick Start
    quickStart: '🚀 クイックスタート',
    downloadJenv: 'JEnv をダウンロード',
    downloadJdk: 'JDK をダウンロード',
    install: 'インストールと使用',
    version: 'バージョン',
    downloading: '読み込み中...',
    download: 'ダウンロード',
    notAvailable: 'お使いのプラットフォームでは利用できません',

    // JDK Selection
    selectDistribution: 'ディストリビューションを選択',
    selectDistributions: 'ディストリビューションを選択',
    selectVersion: 'バージョン',
    selectVersions: 'バージョンを選択',
    yourPlatform: 'プラットフォーム',
    autoDetected: '自動検出',
    jdkHint: '💡 迷っていますか？JDK 11 はほとんどの初心者に最適です',
    distHint: '💡 1つ以上の JDK ディストリビューションを選択してください',
    versionHint: '💡 推奨: JDK 11 または 17',
    selected: '選択済み',
    packages: 'パッケージ',
    downloadSelected: '選択したものをダウンロード',

    // Features
    whyJenv: 'なぜ JEnv なのか？',
    feature1Title: '高速切り替え',
    feature1Desc: 'シンボリックリンクベースで、300ms で JDK を切り替え',
    feature2Title: 'Windows 最適化',
    feature2Desc: 'Windows 10/11 に最適化、権限を自動処理',
    feature3Title: 'モダンな体験',
    feature3Desc: 'カラフルな CLI、ライト/ダークテーマをサポート',

    // FAQ
    faqTitle: 'よくある質問',

    // Footer
    footerText: 'JEnv - Java バージョン管理を簡単に',
    license: 'ライセンス',
    lastUpdated: '最終更新日',

    // Actions
    copy: 'コピー',
    copied: 'クリップボードにコピーされました',
    copyFailed: 'コピーに失敗しました。手動でコピーしてください',
    downloadStarted: 'ダウンロードを開始しました...',

    // Errors
    loading: '読み込み中...',
    errorLoading: 'データの読み込みに失敗しました。ページを更新してください',
    initError: '初期化に失敗しました。ページを更新してください',

    // Mobile
    mobileTitle: 'デスクトップブラウザを使用してください',
    mobileMessage: '最適な体験のために、デスクトップコンピュータでこのページを開いてください'
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
  if (browserLang.toLowerCase().includes('zh')) return 'zh';
  if (browserLang.toLowerCase().includes('ja')) return 'jp';
  return 'en';
}

// Export for use in other scripts
if (typeof module !== 'undefined' && module.exports) {
  module.exports = { t, getLanguage, translations };
}
