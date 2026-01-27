import React, { useState } from 'react';
import {
  SafeAreaView,
  ScrollView,
  StatusBar,
  StyleSheet,
  Text,
  useColorScheme,
  View,
  TouchableOpacity,
  TextInput,
  ActivityIndicator,
} from 'react-native';

import { calculateBoreholeScore, ScoreResult } from './BoreholeBridge';

const App = () => {
  const isDarkMode = useColorScheme() === 'dark';
  const [logs, setLogs] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<ScoreResult | null>(null);

  const handleCalculate = async () => {
    if (!logs.trim()) return;
    setLoading(true);
    const logLines = logs.split('\n').filter(line => line.trim().length > 0);
    const scoreResult = await calculateBoreholeScore(logLines);
    setResult(scoreResult);
    setLoading(false);
  };

  const dashboardColors = {
    good: '#10B981', // Emerald Green
    fair: '#F59E0B', // Sunset Orange
    poor: '#EF4444', // Crimson
    blue: '#2196F3',
  };

  const getHealthData = (score: number) => {
    if (score >= 0.9) return { label: 'Excellent', color: dashboardColors.good };
    if (score >= 0.7) return { label: 'Good', color: dashboardColors.good };
    if (score >= 0.4) return { label: 'Fair', color: dashboardColors.fair };
    return { label: 'Poor', color: dashboardColors.poor };
  };

  const health = getHealthData(result?.score || 0);

  return (
    <SafeAreaView style={[styles.container, isDarkMode ? styles.darkBg : styles.lightBg]}>
      <StatusBar barStyle={isDarkMode ? 'light-content' : 'dark-content'} />
      <ScrollView contentInsetAdjustmentBehavior="automatic">
        <View style={styles.header}>
          <Text style={styles.title}>Borehole Engine</Text>
          <Text style={styles.subtitle}>Fintech Edge Infrastructure</Text>
        </View>

        <View style={styles.card}>
          <Text style={styles.cardTitle}>Financial Data Input</Text>
          <View style={styles.inputContainer}>
            <TextInput
              multiline
              numberOfLines={6}
              style={[styles.input, isDarkMode ? styles.darkInput : styles.lightInput]}
              placeholder="Paste M-Pesa, Airtel, or Bank SMS logs here..."
              placeholderTextColor="#94A3B8"
              value={logs}
              onChangeText={setLogs}
            />
          </View>
          <TouchableOpacity
            style={styles.button}
            onPress={handleCalculate}
            disabled={loading}
          >
            {loading ? (
              <ActivityIndicator color="#fff" />
            ) : (
              <Text style={styles.buttonText}>Calculate Edge Score</Text>
            )}
          </TouchableOpacity>
        </View>

        {result && (
          <View style={styles.resultCard}>
            {/* Speedometer Score Section */}
            <View style={styles.speedometerContainer}>
              <View style={[styles.semiCircle, { borderColor: health.color }]}>
                <Text style={[styles.scoreValue, { color: health.color }]}>
                  {((result.score || 0) * 1000).toFixed(0)}
                </Text>
                <Text style={[styles.healthLabel, { color: health.color }]}>{health.label}</Text>
              </View>
              <Text style={styles.offlineCaption}>Calculated 100% Offline</Text>
            </View>

            {/* The Three Pillars Row */}
            <View style={styles.pillarsRow}>
              <View style={styles.pillar}>
                <Text style={styles.pillarLabel}>Cash In</Text>
                <Text style={styles.pillarValue}>
                  Ksh {result.features?.[0]?.toLocaleString() || '0'}
                </Text>
              </View>
              <View style={[styles.pillar, styles.pillarBorder]}>
                <Text style={styles.pillarLabel}>Cash Out</Text>
                <Text style={styles.pillarValue}>
                  Ksh {result.features?.[1]?.toLocaleString() || '0'}
                </Text>
              </View>
              <View style={styles.pillar}>
                <Text style={styles.pillarLabel}>Debt Level</Text>
                <Text style={[styles.pillarValue, { color: (result.features?.[19] || 0) > 0.3 ? dashboardColors.poor : '#334155' }]}>
                  {((result.features?.[19] || 0) * 100).toFixed(0)}%
                </Text>
              </View>
            </View>

            {result.error && (
              <Text style={styles.errorText}>Error: {result.error}</Text>
            )}
          </View>
        )}

        <View style={styles.footer}>
          <Text style={styles.privacyNote}>🛡️ Privacy First: Your financial logs never leave this device.</Text>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: { flex: 1 },
  lightBg: { backgroundColor: '#F5F7FA' },
  darkBg: { backgroundColor: '#0F172A' },
  header: {
    paddingTop: 30,
    paddingBottom: 20,
    alignItems: 'center',
  },
  title: {
    fontSize: 24,
    fontWeight: '800',
    color: '#0F172A',
    letterSpacing: -0.5,
  },
  subtitle: {
    fontSize: 13,
    color: '#64748B',
    marginTop: 4,
    fontWeight: '500',
  },
  card: {
    marginHorizontal: 20,
    marginBottom: 20,
    padding: 24,
    borderRadius: 16,
    backgroundColor: '#FFFFFF',
    elevation: 4,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.05,
    shadowRadius: 12,
  },
  resultCard: {
    marginHorizontal: 20,
    marginBottom: 30,
    padding: 24,
    borderRadius: 16,
    backgroundColor: '#FFFFFF',
    elevation: 6,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 0.08,
    shadowRadius: 16,
    alignItems: 'center',
  },
  cardTitle: {
    fontSize: 15,
    fontWeight: '700',
    marginBottom: 16,
    color: '#1E293B',
    textTransform: 'uppercase',
    letterSpacing: 0.5,
  },
  inputContainer: {
    borderRadius: 12,
    borderWidth: 1,
    borderColor: '#E2E8F0',
    overflow: 'hidden',
  },
  input: {
    padding: 16,
    fontSize: 15,
    minHeight: 120,
    textAlignVertical: 'top',
  },
  lightInput: { backgroundColor: '#F8FAFC', color: '#1E293B' },
  darkInput: { backgroundColor: '#1E293B', color: '#F1F5F9', borderColor: '#334155' },
  button: {
    backgroundColor: '#2196F3',
    padding: 18,
    borderRadius: 12,
    marginTop: 20,
    alignItems: 'center',
    elevation: 4,
    shadowColor: '#2196F3',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.3,
    shadowRadius: 8,
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '700',
  },
  speedometerContainer: {
    alignItems: 'center',
    marginVertical: 10,
    width: '100%',
  },
  semiCircle: {
    width: 220,
    height: 120,
    borderTopLeftRadius: 110,
    borderTopRightRadius: 110,
    borderWidth: 12,
    borderBottomWidth: 0,
    alignItems: 'center',
    justifyContent: 'flex-end',
    paddingBottom: 5,
  },
  scoreValue: {
    fontSize: 52,
    fontWeight: '900',
    marginBottom: -5,
  },
  healthLabel: {
    fontSize: 18,
    fontWeight: '800',
    textTransform: 'uppercase',
    letterSpacing: 1,
    marginBottom: 5,
  },
  offlineCaption: {
    fontSize: 11,
    color: '#94A3B8',
    marginTop: 15,
    fontWeight: '600',
  },
  pillarsRow: {
    flexDirection: 'row',
    width: '100%',
    marginTop: 30,
    paddingTop: 20,
    borderTopWidth: 1,
    borderTopColor: '#F1F5F9',
  },
  pillar: {
    flex: 1,
    alignItems: 'center',
  },
  pillarBorder: {
    borderLeftWidth: 1,
    borderRightWidth: 1,
    borderColor: '#F1F5F9',
  },
  pillarLabel: {
    fontSize: 11,
    color: '#64748B',
    fontWeight: '700',
    marginBottom: 6,
    textTransform: 'uppercase',
  },
  pillarValue: {
    fontSize: 14,
    fontWeight: '700',
    color: '#334155',
  },
  footer: {
    alignItems: 'center',
    paddingBottom: 40,
  },
  privacyNote: {
    fontSize: 12,
    color: '#94A3B8',
    fontWeight: '500',
  },
  errorText: {
    color: '#EF4444',
    marginTop: 15,
    fontSize: 12,
    textAlign: 'center',
  }
});

export default App;
