// 1. Polyfills must remain at the very top
const TextEncodingPolyfill = require('text-encoding');
Object.assign(global, {
  TextEncoder: TextEncodingPolyfill.TextEncoder,
  TextDecoder: TextEncodingPolyfill.TextDecoder,
});

import React, { useState } from 'react';
import { View, StyleSheet } from 'react-native';

// 2. Import Screens
import WelcomeScreen from './src/screens/WelcomeScreen';
import DashboardScreen from './src/screens/DashboardScreen';

const App = () => {
  // simple state-based navigation
  const [currentScreen, setCurrentScreen] = useState<'welcome' | 'dashboard'>('welcome');

  if (currentScreen === 'welcome') {
    return <WelcomeScreen onStart={() => setCurrentScreen('dashboard')} />;
  }

  return <DashboardScreen />;
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
});

export default App;
