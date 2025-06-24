#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

const BINARY_PATH = path.join(__dirname, '..', 'bin', 'uzp');

function uninstall() {
  try {
    if (fs.existsSync(BINARY_PATH)) {
      fs.unlinkSync(BINARY_PATH);
      console.log('✅ UZP binary removed successfully');
    }
  } catch (error) {
    console.error('❌ Failed to remove UZP binary:', error.message);
  }
}

uninstall(); 