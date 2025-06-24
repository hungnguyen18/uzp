#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const https = require('https');
const { execSync } = require('child_process');

const BINARY_NAME = 'uzp';
const BINARY_PATH = path.join(__dirname, '..', 'bin', BINARY_NAME);
const GITHUB_RELEASES_URL = 'https://api.github.com/repos/hungnguyen18/uzp/releases/latest';

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
    
    https.get(url, (response) => {
      if (response.statusCode === 302 || response.statusCode === 301) {
        // Handle redirect
        return downloadFile(response.headers.location, dest);
      }
      
      if (response.statusCode !== 200) {
        reject(new Error(`HTTP ${response.statusCode}: ${response.statusMessage}`));
        return;
      }
      
      response.pipe(file);
      
      file.on('finish', () => {
        file.close();
        resolve();
      });
      
      file.on('error', (err) => {
        fs.unlink(dest, () => {});
        reject(err);
      });
    }).on('error', (err) => {
      reject(err);
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
          const release = JSON.parse(data);
          resolve(release);
        } catch (err) {
          reject(err);
        }
      });
    }).on('error', (err) => {
      reject(err);
    });
  });
}

async function install() {
  try {
    console.log('📦 Installing UZP...');
    
    const binaryName = getBinaryName();
    console.log(`🔍 Looking for binary: ${binaryName}`);
    
    // Get latest release info
    const release = await getLatestRelease();
    
    // Find the asset for current platform
    const asset = release.assets.find(asset => asset.name === binaryName);
    
    if (!asset) {
      throw new Error(`Binary not found for platform: ${binaryName}`);
    }
    
    console.log(`⬇️  Downloading from: ${asset.browser_download_url}`);
    
    // Download binary
    await downloadFile(asset.browser_download_url, BINARY_PATH);
    
    // Make executable on Unix-like systems
    if (process.platform !== 'win32') {
      fs.chmodSync(BINARY_PATH, '755');
    }
    
    console.log('✅ UZP installed successfully!');
    console.log('');
    console.log('🚀 Get started:');
    console.log('   uzp init');
    console.log('   uzp add');
    console.log('   uzp --help');
    
  } catch (error) {
    console.error('❌ Installation failed:', error.message);
    console.log('');
    console.log('🔧 Manual installation:');
    console.log('   git clone https://github.com/hungnguyen18/uzp.git');
    console.log('   cd uzp');
    console.log('   go build -o uzp');
    process.exit(1);
  }
}

// Run installation
install(); 