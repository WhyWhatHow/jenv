# JEnv Landing Page

A modern, responsive landing page for JEnv - the fast and easy Java version manager.

## Features

- 🎯 **Platform Detection**: Automatically detects user's OS and architecture (Windows, Linux, macOS with x64/arm64)
- 📦 **One-Click Downloads**: Direct download links for JEnv and JDK distributions
- 🌍 **Internationalization**: Supports English and Chinese with automatic language detection
- 🎨 **Modern UI**: Dark theme with smooth animations and responsive design
- 📱 **Mobile Friendly**: Optimized for desktop with mobile fallback page
- 🤖 **Auto-Updates**: GitHub Actions automatically fetch latest JDK links every 6 hours
- ⚡ **Zero Runtime Cost**: Pure static site with pre-generated data

## Quick Start

### Prerequisites

- Node.js 18+ (for development only)
- Modern web browser

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/WhyWhatHow/jenv.git
cd jenv/landing-page
```

2. Serve the static files (choose one):
```bash
# Using Python
python -m http.server 8000

# Using Node.js
npx serve .

# Using PHP
php -S localhost:8000
```

3. Open your browser:
```
http://localhost:8000
```

### Manual Data Update

To manually fetch JDK links:

```bash
cd scripts
node fetch-jdk-links.js
```

This will update `data/jdk.json` with the latest JEnv releases and JDK download links.

## Project Structure

```
landing-page/
├── index.html              # Main HTML file
├── css/
│   └── style.css          # Styles with responsive design
├── js/
│   ├── app.js             # Main application logic
│   └── i18n.js            # Internationalization support
├── data/
│   └── jdk.json           # JDK and JEnv download data (auto-generated)
├── scripts/
│   └── fetch-jdk-links.js # Script to fetch JDK links
├── .github/
│   └── workflows/
│       └── fetch-jdk-links.yml  # Auto-update workflow
└── README.md              # This file
```

## Deployment

### GitHub Pages

1. Enable GitHub Pages in your repository settings:
   - Go to **Settings** → **Pages**
   - Source: **Deploy from a branch**
   - Branch: `main` (or your preferred branch)
   - Folder: `/landing-page`

2. The site will be available at:
   ```
   https://yourusername.github.io/jenv/
   ```

### Custom Domain (Optional)

1. Add a `CNAME` file:
```bash
echo "jenv.yourdomain.com" > CNAME
```

2. Configure DNS:
   - Add a CNAME record pointing to `yourusername.github.io`

3. Enable HTTPS in repository settings

### Netlify

1. Connect your GitHub repository
2. Configure build settings:
   - Base directory: `landing-page`
   - Publish directory: `landing-page`
   - Build command: (leave empty)

3. Deploy!

### Vercel

1. Import your GitHub repository
2. Configure project:
   - Root directory: `landing-page`
   - Framework preset: Other
   - Build command: (leave empty)
   - Output directory: `.`

3. Deploy!

## Automated Maintenance

### JDK Link Updates

The GitHub Actions workflow (`.github/workflows/fetch-jdk-links.yml`) runs automatically:
- **Schedule**: Every 6 hours
- **Manual Trigger**: Can be triggered manually from Actions tab

The workflow:
1. Fetches latest JEnv releases from GitHub API
2. Fetches JDK download links from Adoptium API
3. Updates `data/jdk.json`
4. Commits and pushes changes

### Adding New JDK Distributions

To add support for more JDK distributions (Zulu, Corretto, etc.):

1. Edit `scripts/fetch-jdk-links.js`
2. Add new distribution fetcher function:
```javascript
async function fetchZuluJDK(version, platform) {
  // Implement Zulu API fetching
}
```

3. Update the main function to include new distribution:
```javascript
const zulu = await fetchZuluDistribution();
const output = {
  jdk: {
    distributions: {
      temurin,
      zulu  // Add new distribution
    }
  }
};
```

4. Test locally:
```bash
node scripts/fetch-jdk-links.js
```

## Technology Stack

- **HTML5**: Semantic markup with accessibility features
- **CSS3**: Modern CSS with custom properties, Grid, and Flexbox
- **Vanilla JavaScript**: No framework dependencies
- **GitHub Actions**: Automated data fetching
- **Node.js**: Script runtime for data fetching

## Browser Support

- Chrome/Edge: Latest 2 versions
- Firefox: Latest 2 versions
- Safari: Latest 2 versions

## Contributing

1. Fork the repository
2. Create your feature branch: `git checkout -b feat/amazing-feature`
3. Commit your changes: `git commit -m "feat: add amazing feature"`
4. Push to the branch: `git push origin feat/amazing-feature`
5. Open a Pull Request

## License

This project is part of JEnv and follows the same license.

## Links

- [JEnv GitHub Repository](https://github.com/WhyWhatHow/jenv)
- [Report Issues](https://github.com/WhyWhatHow/jenv/issues)
- [JEnv Documentation](https://github.com/WhyWhatHow/jenv/blob/main/README.md)
