# JDK 下拉选择功能

## 功能概述

Landing Page 现在提供简洁的下拉菜单来选择 JDK 发行版和版本，一键下载所需的 JDK。

## 实现的功能

### 1. 多供应商支持

- **集成 Foojay DiscoAPI**: 使用 foojay.io API 获取多个 JDK 发行版的数据
- **当前支持的发行版**:
  - Eclipse Temurin (推荐) ⭐
  - Azul Zulu
  - Amazon Corretto

- **可扩展支持**（取消注释即可启用）:
  - BellSoft Liberica
  - Microsoft Build of OpenJDK
  - Oracle OpenJDK
  - GraalVM CE 17/21
  - SapMachine
  - Alibaba Dragonwell

### 2. 完整版本支持

- **支持的 JDK 版本**: 8, 11, 17, 21, 25
- **推荐版本**: 11, 17 (带星标 ⭐)
- **平台支持**: Windows x64, Linux x64, Linux ARM64, macOS Intel, macOS Apple Silicon
- **自动更新**: 可随时添加新版本，只需修改 `JDK_VERSIONS` 数组

### 3. 优化的 UI 设计

#### 下拉菜单
- 简洁美观的下拉选择器
- 推荐项目带星标标识
- 自动选择推荐版本
- 悬停和聚焦状态动画

#### 布局优化
- 两栏等高卡片布局
- JEnv 和 JDK 下载框统一高度（350px）
- Flexbox 自适应内容分布
- 响应式设计，移动端自动变为单列

#### 多语言支持
- 完整的中英文界面
- 所有文案已翻译

### 4. 智能下载功能

- 根据用户选择的发行版和版本自动匹配平台
- 下载按钮显示完整信息（发行版名称 + 版本 + 文件大小）
- 自动检测是否有可用下载
- 平台不支持时自动禁用按钮

## 文件变更

### 修改的文件

1. **index.html**
   - 简化为下拉菜单选择器
   - 优化卡片结构和布局

2. **css/style.css** (+205 行)
   - 两栏等高布局
   - 下拉菜单样式（含自定义箭头）
   - Flexbox 内容分布
   - 响应式断点优化

3. **js/app.js** (+116 行改动)
   - 简化为下拉选择逻辑
   - 动态填充发行版选项
   - 动态填充版本选项
   - 实时更新下载按钮

4. **js/i18n.js** (+16 行)
   - 添加新的翻译键
   - 支持英文和中文

5. **scripts/fetch-jdk-links.js** (+118 行改动)
   - 使用 Foojay DiscoAPI
   - 支持 JDK 25
   - 支持 10+ 个发行版
   - 优化 API 调用

6. **data/jdk.json** (+553 行)
   - 包含 3 个发行版数据
   - 5 个 JDK 版本 (8, 11, 17, 21, 25)
   - 5 个平台支持

## 数据格式

### 生成的 JSON 格式

```json
{
  "lastUpdated": "2025-11-25T14:17:27.367Z",
  "jenv": {
    "version": "0.6.7",
    "platforms": { ... }
  },
  "jdk": {
    "versions": [8, 11, 17, 21, 22, 23],
    "recommended": [11, 17],
    "distributions": {
      "temurin": {
        "name": "Eclipse Temurin",
        "description": "Most popular open-source JDK",
        "recommended": true,
        "versions": {
          "11": {
            "windows-x64": {
              "url": "https://...",
              "size": "190.1 MB",
              "sha256": "...",
              "javaVersion": "11.0.29+7",
              "distribution": "temurin"
            }
          }
        }
      }
    }
  }
}
```

## 使用方法

### 本地开发

1. 生成 JDK 数据:
```bash
cd landing-page/scripts
node fetch-jdk-links.js
```

2. 启动本地服务器:
```bash
cd landing-page
python3 -m http.server 8000
```

3. 访问: http://localhost:8000

### 配置发行版

编辑 `scripts/fetch-jdk-links.js` 中的 `DISTRIBUTIONS` 数组:

```javascript
const DISTRIBUTIONS = [
  { id: 'temurin', name: 'Eclipse Temurin', desc: 'Most popular open-source JDK', recommended: true },
  { id: 'zulu', name: 'Azul Zulu', desc: 'Enterprise-ready OpenJDK' },
  // 取消注释下面的行来添加更多发行版
  // { id: 'liberica', name: 'BellSoft Liberica', desc: 'Flexible OpenJDK builds' },
];
```

### 配置版本

修改 `JDK_VERSIONS` 数组:

```javascript
const JDK_VERSIONS = [8, 11, 17, 21, 25]; // 添加新版本很简单
```

## UI 预览

### 布局特点

- **等高卡片**: 两个下载框高度统一（350px），视觉更协调
- **下拉选择**: 简洁的下拉菜单替代复杂的复选框
- **推荐标识**: 推荐的发行版和版本带星标 ⭐
- **自适应**: 内容使用 Flexbox 自动分布，按钮始终在底部

### 响应式断点

- **桌面** (> 768px): 两栏布局
- **移动端** (≤ 768px): 单栏布局，自动调整

## 性能优化

1. **API 调用优化**
   - 每次请求之间延迟 100ms
   - 避免 API 限流
   - 智能重试机制

2. **UI 性能**
   - CSS 变量复用
   - 最小化 DOM 操作
   - 使用原生 select 元素

3. **数据优化**
   - 按需加载发行版数据
   - 静态 JSON 缓存

## 浏览器兼容性

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+
- 移动端浏览器

## 与旧版本的对比

### 旧版本（多选复选框）
- ❌ 界面复杂，需要多次点击
- ❌ 占用空间大
- ❌ 移动端体验不佳
- ✅ 可以同时下载多个包

### 新版本（下拉菜单）
- ✅ 界面简洁，操作直观
- ✅ 节省空间，布局协调
- ✅ 移动端体验优秀
- ✅ 符合传统下载页面习惯
- ✅ 加载速度更快

## API 来源

- **JEnv 版本**: GitHub Releases API
- **JDK 包**: Foojay DiscoAPI (https://api.foojay.io/disco/v3.0/)
  - 支持 20+ 个 JDK 发行版
  - 实时更新的下载链接
  - 包含 SHA256 校验和

## 未来改进

1. ✅ 添加 JDK 25 支持
2. ✅ 优化布局和 UI
3. ⬜ 添加下载统计
4. ⬜ 支持用户自定义选择保存
5. ⬜ 添加 JDK 版本对比功能
6. ⬜ 集成安装脚本生成器
7. ⬜ 添加自动检测最新 JDK 版本

## 常见问题

### Q: 如何添加新的 JDK 版本？
A: 编辑 `scripts/fetch-jdk-links.js`，在 `JDK_VERSIONS` 数组中添加版本号，然后重新运行脚本。

### Q: 为什么改为下拉菜单？
A: 下拉菜单更简洁，符合传统下载页面的使用习惯，布局更协调，移动端体验更好。

### Q: 如何添加更多发行版？
A: 在 `scripts/fetch-jdk-links.js` 的 `DISTRIBUTIONS` 数组中取消注释相应的发行版，或添加新的发行版配置。

### Q: 数据多久更新一次？
A: 可以通过 GitHub Actions 配置定时任务自动更新，建议每天或每周更新一次。

## 相关链接

- [Foojay DiscoAPI 文档](https://github.com/foojayio/discoapi)
- [JEnv GitHub](https://github.com/WhyWhatHow/jenv)
- [Foojay API Endpoint](https://api.foojay.io/disco/v3.0/)

---

**最后更新**: 2025-11-26
**版本**: 2.0.0
**作者**: Claude Code
**变更**: 从多选复选框改为下拉菜单，添加 JDK 25 支持，优化布局
