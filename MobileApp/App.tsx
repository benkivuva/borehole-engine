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

  const getScoreColor = (score: number) => {
    if (score > 0.7) return '#4CAF50';
    if (score > 0.4) return '#FFC107';
    return '#F44336';
  };

  return (
    <SafeAreaView style={[styles.container, isDarkMode ? styles.darkBg : styles.lightBg]}>
      <StatusBar barStyle={isDarkMode ? 'light-content' : 'dark-content'} />
      <ScrollView contentInsetAdjustmentBehavior="automatic">
        <View style={styles.header}>
          <Text style={styles.title}>Borehole Engine</Text>
          <Text style={styles.subtitle}>Edge Credit Scoring Infrastructure</Text>
        </View>

        <View style={styles.card}>
          <Text style={styles.cardTitle}>SMS Logs Input</Text>
          <TextInput
            multiline
            numberOfLines={10}
            style={[styles.input, isDarkMode ? styles.darkInput : styles.lightInput]}
            placeholder="Paste M-Pesa, Airtel, or Bank SMS logs here..."
            placeholderTextColor="#888"
            value={logs}
            onChangeText={setLogs}
          />
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
            <Text style={styles.cardTitle}>Scoring Result</Text>
            <View style={styles.scoreContainer}>
              <Text style={[styles.scoreValue, { color: getScoreColor(result.score || 0) }]}>
                {((result.score || 0) * 1000).toFixed(0)}
              </Text>
              <Text style={styles.scoreLabel}>Borehole Index</Text>
            </View>

            <View style={styles.statsRow}>
              <View style={styles.statBox}>
                <Text style={styles.statValue}>{result.txn_count || 0}</Text>
                <Text style={styles.statLabel}>Transactions</Text>
              </View>
              <View style={styles.statBox}>
                <Text style={styles.statValue}>Verified</Text>
                <Text style={styles.statLabel}>Status</Text>
              </View>
            </View>

            {(result?.features?.length || 0) > 0 && (
              <View style={styles.featureList}>
                <Text style={styles.featureTitle}>Feature Vector (Selected High-Impact):</Text>
                <View style={styles.featureRow}>
                  <Text style={styles.featureText}>Hustler Balance: {result.features?.[11]?.toFixed(2) || '0.00'}</Text>
                  <Text style={styles.featureText}>Okoa Reliance: {result.features?.[12]?.toFixed(2) || '0.00'}</Text>
                </View>
              </View>
            )}

            {result.error && (
              <Text style={styles.errorText}>Error: {result.error}</Text>
            )}
          </View>
        )}
      </ScrollView>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  lightBg: { backgroundColor: '#F5F7FA' },
  darkBg: { backgroundColor: '#121212' },
  header: {
    padding: 30,
    alignItems: 'center',
  },
  title: {
    fontSize: 28,
    fontWeight: '800',
    color: '#2196F3',
    letterSpacing: 1,
  },
  subtitle: {
    fontSize: 14,
    color: '#666',
    marginTop: 5,
  },
  card: {
    margin: 20,
    padding: 20,
    borderRadius: 16,
    backgroundColor: '#fff',
    elevation: 4,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 8,
  },
  resultCard: {
    margin: 20,
    marginTop: 0,
    padding: 25,
    borderRadius: 16,
    backgroundColor: '#fff',
    elevation: 4,
    alignItems: 'center',
  },
  cardTitle: {
    fontSize: 16,
    fontWeight: '700',
    marginBottom: 15,
    color: '#333',
  },
  input: {
    borderRadius: 12,
    padding: 15,
    fontSize: 14,
    minHeight: 150,
    textAlignVertical: 'top',
    borderWidth: 1,
    borderColor: '#E0E0E0',
  },
  lightInput: { backgroundColor: '#FAFAFA', color: '#333' },
  darkInput: { backgroundColor: '#1E1E1E', color: '#EEE', borderColor: '#333' },
  button: {
    backgroundColor: '#2196F3',
    padding: 18,
    borderRadius: 12,
    marginTop: 20,
    alignItems: 'center',
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '700',
  },
  scoreContainer: {
    alignItems: 'center',
    marginVertical: 20,
  },
  scoreValue: {
    fontSize: 64,
    fontWeight: '900',
  },
  scoreLabel: {
    fontSize: 16,
    color: '#888',
    fontWeight: '600',
  },
  statsRow: {
    flexDirection: 'row',
    width: '100%',
    justifyContent: 'space-around',
    borderTopWidth: 1,
    borderTopColor: '#EEE',
    paddingTop: 20,
  },
  statBox: {
    alignItems: 'center',
  },
  statValue: {
    fontSize: 20,
    fontWeight: '700',
    color: '#333',
  },
  statLabel: {
    fontSize: 12,
    color: '#999',
    marginTop: 4,
  },
  featureList: {
    marginTop: 20,
    width: '100%',
    padding: 15,
    backgroundColor: '#F8F9FA',
    borderRadius: 8,
  },
  featureTitle: {
    fontSize: 12,
    fontWeight: '700',
    color: '#666',
    marginBottom: 8,
  },
  featureRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  featureText: {
    fontSize: 11,
    color: '#444',
    fontFamily: 'monospace',
  },
  errorText: {
    color: '#F44336',
    marginTop: 10,
    fontSize: 12,
  }
});

export default App;
