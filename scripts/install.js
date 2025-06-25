#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const https = require('https');
const os = require('os');
const { execSync } = require('child_process');

const BINARY_NAME = 'uzp';
const BINARY_PATH = path.join(__dirname, '..', 'bin', BINARY_NAME);
const GITHUB_RELEASES_URL = 'https://api.github.com/repos/hungnguyen18/uzp-cli/releases/latest';

// Cache directory in user's home
const CACHE_DIR = path.join(os.homedir(), '.uzp-cache');
const VERSION_FILE = path.join(CACHE_DIR, 'version.txt');

function ensureCacheDir() {
  if (!fs.existsSync(CACHE_DIR)) {
    fs.mkdirSync(CACHE_DIR, { recursive: true });
  }
}

function ensureBinDir() {
  const binDir = path.dirname(BINARY_PATH);
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }
}

function getCachedBinaryPath(binaryName, version) {
  return path.join(CACHE_DIR, `${binaryName}-${version}`);
}

function isBinaryCached(binaryName, version) {
  const cachedPath = getCachedBinaryPath(binaryName, version);
  return fs.existsSync(cachedPath);
}

function copyFromCache(binaryName, version, dest) {
  const cachedPath = getCachedBinaryPath(binaryName, version);
  fs.copyFileSync(cachedPath, dest);
}

function cacheDownloadedBinary(binaryName, version, source) {
  ensureCacheDir();
  const cachedPath = getCachedBinaryPath(binaryName, version);
  fs.copyFileSync(source, cachedPath);
  // Store version without 'v' prefix for consistency
  const versionToStore = version.startsWith('v') ? version.substring(1) : version;
  fs.writeFileSync(VERSION_FILE, versionToStore, 'utf8');
}

function getPlatform() {
  switch (process.platform) {
    case 'darwin':
      return 'darwin';
    case 'linux':
      return 'linux';
    case 'win32':
      return 'windows';
    default:
      throw new Error(`Unsupported platform: ${process.platform}`);
  }
}

function getArch() {
  switch (process.arch) {
    case 'x64':
      return 'amd64';
    case 'arm64':
      return 'arm64';
    default:
      throw new Error(`Unsupported architecture: ${process.arch}`);
  }
}

function getBinaryName() {
  const platform = getPlatform();
  const arch = getArch();
  const ext = platform === 'windows' ? '.exe' : '';
  return `uzp-${platform}-${arch}${ext}`;
}

function downloadFile(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    let resolved = false;
    
    const cleanup = () => {
      if (!resolved) {
        resolved = true;
        file.destroy();
      }
    };
    
    https.get(url, (response) => {
      if (response.statusCode === 302 || response.statusCode === 301) {
        file.destroy();
        return downloadFile(response.headers.location, dest).then(resolve).catch(reject);
      }
      
      if (response.statusCode !== 200) {
        cleanup();
        reject(new Error(`HTTP ${response.statusCode}: ${response.statusMessage}`));
        return;
      }
      
      const totalSize = parseInt(response.headers['content-length'], 10);
      let downloadedSize = 0;
      let lastPercentage = 0;
      
      if (totalSize) {
        response.on('data', (chunk) => {
          downloadedSize += chunk.length;
          const percentage = Math.floor((downloadedSize / totalSize) * 100);
          
          if (percentage > lastPercentage && percentage % 20 === 0) {
            process.stdout.write(`\r⬇️  Downloading... ${percentage}%`);
            lastPercentage = percentage;
          }
        });
      }
      
      response.pipe(file);
      
      file.on('finish', () => {
        if (!resolved) {
          resolved = true;
          if (totalSize) {
            process.stdout.write('\r⬇️  Download complete!     \n');
          }
          resolve();
        }
      });
      
      file.on('error', (err) => {
        cleanup();
        fs.unlink(dest, () => {});
        if (!resolved) {
          resolved = true;
          reject(err);
        }
      });
      
      response.on('error', (err) => {
        cleanup();
        if (!resolved) {
          resolved = true;
          reject(err);
        }
      });
    }).on('error', (err) => {
      cleanup();
      if (!resolved) {
        resolved = true;
        reject(err);
      }
    });
  });
}

function getRequestedVersion() {
  // Try different package manager environment variables
  const sources = [
    process.env.npm_config_wanted_version,     // npm install uzp-cli@1.0.6
    process.env.npm_package_version,           // from package.json
    process.env.YARN_WANTED_VERSION,           // yarn add global uzp-cli@1.0.6
    process.env.PNPM_WANTED_VERSION,           // pnpm add -g uzp-cli@1.0.6
    process.env.BUN_WANTED_VERSION,            // bun add -g uzp-cli@1.0.6
    process.argv.find(arg => arg.includes('@')), // command line uzp-cli@1.0.6
  ];
  
  for (const source of sources) {
    if (source) {
      // Extract version from uzp-cli@1.0.6 format
      const match = source.match(/@(.+)$/);
      if (match) {
        return match[1];
      }
      // If it's just a version number
      if (source.match(/^\d+\.\d+\.\d+/)) {
        return source;
      }
    }
  }
  
  return null;
}

async function getSpecificRelease(version) {
  const tagName = version.startsWith('v') ? version : `v${version}`;
  const url = `https://api.github.com/repos/hungnguyen18/uzp-cli/releases/tags/${tagName}`;
  
  return new Promise((resolve, reject) => {
    https.get(url, {
      headers: {
        'User-Agent': 'uzp-npm-installer'
      }
    }, (response) => {
      let data = '';
      
      response.on('data', (chunk) => {
        data += chunk;
      });
      
      response.on('end', () => {
        try {
          if (response.statusCode === 404) {
            reject(new Error(`Version ${version} not found. Check available versions at: https://github.com/hungnguyen18/uzp-cli/releases`));
            return;
          }
          
          if (response.statusCode !== 200) {
            reject(new Error(`GitHub API error: ${response.statusCode} ${response.statusMessage}`));
            return;
          }
          
          const release = JSON.parse(data);
          
          // Check if release has assets
          if (!release.assets || !Array.isArray(release.assets) || release.assets.length === 0) {
            reject(new Error(`No release assets found for version ${version}.`));
            return;
          }
          
          resolve(release);
        } catch (err) {
          reject(new Error(`Failed to parse GitHub release data: ${err.message}`));
        }
      });
    }).on('error', (err) => {
      reject(new Error(`Failed to fetch GitHub release: ${err.message}`));
    });
  });
}

async function getLatestRelease() {
  return new Promise((resolve, reject) => {
    https.get(GITHUB_RELEASES_URL, {
      headers: {
        'User-Agent': 'uzp-npm-installer'
      }
    }, (response) => {
      let data = '';
      
      response.on('data', (chunk) => {
        data += chunk;
      });
      
      response.on('end', () => {
        try {
          if (response.statusCode === 404) {
            reject(new Error('No GitHub releases found. Please create a release first.'));
            return;
          }
          
          if (response.statusCode !== 200) {
            reject(new Error(`GitHub API error: ${response.statusCode} ${response.statusMessage}`));
            return;
          }
          
          const release = JSON.parse(data);
          
          // Check if release has assets
          if (!release.assets || !Array.isArray(release.assets) || release.assets.length === 0) {
            reject(new Error('No release assets found. Please upload binaries to the GitHub release.'));
            return;
          }
          
          resolve(release);
        } catch (err) {
          reject(new Error(`Failed to parse GitHub release data: ${err.message}`));
        }
      });
    }).on('error', (err) => {
      reject(new Error(`Failed to fetch GitHub release: ${err.message}`));
    });
  });
}

function checkExistingInstallation() {
  try {
    // Check if binary already exists
    if (fs.existsSync(BINARY_PATH)) {
      console.log('⚠️  UZP binary already exists at:', BINARY_PATH);
      
      // Try to get version from cache or binary
      try {
        // First try to get version from cache file
        if (fs.existsSync(VERSION_FILE)) {
          const cachedVersion = fs.readFileSync(VERSION_FILE, 'utf8').trim();
          console.log(`   Current version: v${cachedVersion} (from cache)`);
          return { exists: true, version: `v${cachedVersion}` };
        }
        
        // Fallback: try binary help command (uzp doesn't have --version)
        const helpOutput = execSync(`"${BINARY_PATH}" --help`, { encoding: 'utf8', stdio: 'pipe' });
        console.log('   Current version: unknown (no version info available)');
        return { exists: true, version: 'unknown' };
      } catch (error) {
        console.log('   Unable to determine current version');
        return { exists: true, version: 'unknown' };
      }
    }
    
    return { exists: false };
  } catch (error) {
    return { exists: false };
  }
}

function suggestReinstall() {
  console.log('');
  console.log('🔧 To fix EEXIST errors, try one of these solutions:');
  console.log('');
  console.log('   Option 1: Use our reinstall tool (recommended)');
  console.log('   npx uzp-reinstall');
  console.log('');
  console.log('   Option 2: Cleanup then reinstall');
  console.log('   npx uzp-cleanup');
  console.log('   npm install -g uzp-cli');
  console.log('');
  console.log('   Option 3: Manual cleanup then reinstall');
  console.log('   npm uninstall -g uzp-cli');
  console.log('   npm cache clean --force');
  console.log('   npm install -g uzp-cli');
  console.log('');
  console.log('   Option 4: Force overwrite (use with caution)');
  console.log('   npm install -g uzp-cli --force');
}

async function install() {
  try {
    console.log('📦 Installing UZP...');
    
    const binaryName = getBinaryName();
    console.log(`🔍 Looking for binary: ${binaryName}`);
    
    // Check for existing installation
    const existingInstall = checkExistingInstallation();
    
    // Check if specific version was requested
    const requestedVersion = getRequestedVersion();
    console.log(`🎯 Target: ${requestedVersion ? `v${requestedVersion}` : 'latest version'}`);
    
    // Get release info (specific version or latest)
    let release;
    if (requestedVersion && requestedVersion !== 'latest') {
      console.log(`📦 Fetching specific version: ${requestedVersion}`);
      release = await getSpecificRelease(requestedVersion);
    } else {
      console.log('📦 Fetching latest version');
      release = await getLatestRelease();
    }
    
    const version = release.tag_name;
    console.log(`📋 Selected version: ${version}`);
    
    // Handle existing installation with override behavior
    if (existingInstall.exists) {
      const isSameVersion = existingInstall.version === version || existingInstall.version.includes(version.replace('v', ''));
      
      if (requestedVersion) {
        // Specific version requested - always override
        console.log(`🔄 Override requested: Installing ${version} (current: ${existingInstall.version})`);
        console.log('🗑️ Removing existing installation...');
        
        try {
          if (fs.existsSync(BINARY_PATH)) {
            fs.unlinkSync(BINARY_PATH);
            console.log('   ✅ Existing binary removed');
          }
        } catch (error) {
          console.log('   ⚠️ Could not remove existing binary, continuing...');
        }
      } else if (isSameVersion) {
        // Latest requested and same version - skip unless forced
        console.log(`✅ UZP ${version} is already installed and up to date!`);
        console.log('');
        console.log('💡 To force reinstall: npx uzp-reinstall');
        console.log('💡 To install specific version: npm install -g uzp-cli@1.0.6');
        console.log('🚀 Ready to use: uzp --help');
        return;
      } else {
        // Latest requested and different version - update
        console.log(`🔄 Updating from ${existingInstall.version} to ${version}`);
        console.log('🗑️ Removing existing installation...');
        
        try {
          if (fs.existsSync(BINARY_PATH)) {
            fs.unlinkSync(BINARY_PATH);
            console.log('   ✅ Existing binary removed');
          }
        } catch (error) {
          console.log('   ⚠️ Could not remove existing binary, continuing...');
        }
      }
    }
    
    // Check if binary is already cached
    if (isBinaryCached(binaryName, version)) {
      console.log(`⚡ Using cached binary for version ${version}`);
      
      try {
        // Ensure bin directory exists before copying
        ensureBinDir();
        copyFromCache(binaryName, version, BINARY_PATH);
        
        // Make executable on Unix-like systems
        if (process.platform !== 'win32') {
          fs.chmodSync(BINARY_PATH, '755');
        }
        
        // Create symlink manually for cached installation too
        try {
          const npmBin = execSync('npm bin -g', { encoding: 'utf8' }).trim();
          const symlinkPath = path.join(npmBin, BINARY_NAME);
          
          // Remove existing symlink if it exists
          if (fs.existsSync(symlinkPath)) {
            fs.unlinkSync(symlinkPath);
          }
          
          // Create new symlink
          fs.symlinkSync(BINARY_PATH, symlinkPath);
          console.log(`🔗 Created symlink: ${symlinkPath} -> ${BINARY_PATH}`);
        } catch (error) {
          console.log(`⚠️  Could not create symlink (${error.message})`);
        }

        console.log('✅ UZP installed successfully from cache!');
        console.log('');
        console.log('🚀 Get started:');
        console.log('   uzp init');
        console.log('   uzp add');
        console.log('   uzp --help');
        return;
      } catch (error) {
        if (error.code === 'EEXIST' || error.message.includes('file already exists')) {
          console.log('❌ File already exists and cannot be overwritten');
          suggestReinstall();
          process.exit(1);
        }
        throw error;
      }
    }
    
    // Find the asset for current platform
    const asset = release.assets.find(asset => asset.name === binaryName);
    
    if (!asset) {
      const availableAssets = release.assets.map(a => a.name).join(', ');
      throw new Error(`Binary not found for platform: ${binaryName}\nAvailable assets: ${availableAssets}`);
    }
    
    console.log(`⬇️  Downloading ${(asset.size / 1024 / 1024).toFixed(1)}MB from GitHub...`);
    
    // Download binary
    try {
      // Ensure bin directory exists before downloading
      ensureBinDir();
      await downloadFile(asset.browser_download_url, BINARY_PATH);
      
      // Cache the downloaded binary
      cacheDownloadedBinary(binaryName, version, BINARY_PATH);
      
      // Make executable on Unix-like systems
      if (process.platform !== 'win32') {
        fs.chmodSync(BINARY_PATH, '755');
      }
    } catch (error) {
      if (error.code === 'EEXIST' || error.message.includes('file already exists')) {
        console.log('❌ Cannot download - file already exists');
        suggestReinstall();
        process.exit(1);
      }
      throw error;
    }
    
    // Create symlink manually (NPM doesn't create it when binary is downloaded in postinstall)
    try {
      const npmBin = execSync('npm bin -g', { encoding: 'utf8' }).trim();
      const symlinkPath = path.join(npmBin, BINARY_NAME);
      
      // Remove existing symlink if it exists
      if (fs.existsSync(symlinkPath)) {
        fs.unlinkSync(symlinkPath);
      }
      
      // Create new symlink
      fs.symlinkSync(BINARY_PATH, symlinkPath);
      console.log(`🔗 Created symlink: ${symlinkPath} -> ${BINARY_PATH}`);
    } catch (error) {
      console.log(`⚠️  Could not create symlink (${error.message})`);
      console.log('   The binary is installed but may not be in PATH');
      console.log(`   Binary location: ${BINARY_PATH}`);
    }

    console.log('✅ UZP installed successfully!');
    console.log('💾 Binary cached for future installations');
    console.log('');
    console.log('🚀 Get started:');
    console.log('   uzp init');
    console.log('   uzp add');
    console.log('   uzp --help');
    
  } catch (error) {
    console.error('❌ Installation failed:', error.message);
    console.log('');
    
    if (error.message.includes('No GitHub releases found') || error.message.includes('No release assets found')) {
      console.log('📋 This package requires a GitHub release with pre-built binaries.');
      console.log('   The maintainer needs to create a release at:');
      console.log('   https://github.com/hungnguyen18/uzp-cli/releases');
      console.log('');
    }
    
    console.log('🔧 Manual installation:');
         console.log('   git clone https://github.com/hungnguyen18/uzp-cli.git');
     console.log('   cd uzp-cli');
    console.log('   go build -o uzp');
    console.log('   sudo mv uzp /usr/local/bin/  # Optional: make globally available');
    process.exit(1);
  }
}

// Run installation
install(); 