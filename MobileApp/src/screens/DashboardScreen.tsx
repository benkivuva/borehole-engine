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
    PermissionsAndroid,
    Alert,
    Modal,
    FlatList,
    Vibration,
    Share,
} from 'react-native';
import SmsAndroid from 'react-native-get-sms-android';
import QRCode from 'react-native-qrcode-svg';

// Adjust relative paths
import { Database, AuditLog } from '../storage/Database';
import { calculateBoreholeScore, ScoreResult, generateSignedScore, SignedCertificate } from '../../BoreholeBridge';

// Mock Data for Quick Testing
const MOCK_PERSONAS = {
    high_income: "RC9999ZZ Confirmed. You have received Ksh75,000.00 from ELON MUSK on 28/1/26 at 1:00 PM.\nUB1234AB Confirmed. Ksh5,000.00 sent to MOM on 29/1/26.",
    average: "RC8888AA Confirmed. You have received Ksh15,000.00 from SALARY on 28/1/26.\nFuliza M-PESA. You have borrowed Ksh500.00.",
    risky: "RC7777BB Confirmed. Ksh5,000.00 sent to BETIKA on 28/1/26.\nFuliza M-PESA. You have borrowed Ksh2,000.00.\nOkoa Jahazi: You have a debt of Ksh50.",
};

const DashboardScreen = () => {
    const isDarkMode = useColorScheme() === 'dark';
    const [logs, setLogs] = useState('');
    const [loading, setLoading] = useState(false);
    const [scanning, setScanning] = useState(false);
    const [result, setResult] = useState<ScoreResult | null>(null);
    const [history, setHistory] = useState<AuditLog[]>([]);
    const [showHistory, setShowHistory] = useState(false);
    const [showVerify, setShowVerify] = useState(false);
    const [cert, setCert] = useState<SignedCertificate | null>(null);

    const loadHistory = async () => {
        const data = await Database.getHistory();
        setHistory(data);
        setShowHistory(true);
    };

    const nukeHistory = async () => {
        await Database.nukeData();
        setHistory([]);
        setResult(null);
        setLogs('');
        setCert(null);
        setShowHistory(false);
        Alert.alert('System Reset', 'Local database and dashboard cleared.');
    };

    const handleVerify = async () => {
        if (!result || !result.score) return;
        Vibration.vibrate(10);
        setLoading(true);
        const certificate = await generateSignedScore(result.score);
        setCert(certificate);
        setLoading(false);
        Vibration.vibrate(50);
        setShowVerify(true);
    };

    const handleShare = async () => {
        if (!cert) return;
        try {
            await Share.share({
                message: `🔒 Verified Borehole Score: ${((result?.score || 0) * 1000).toFixed(0)}\n\nSignature:\n${cert.signature}\n\nVerify at: https://borehole.fi`,
            });
        } catch (error) {
            console.error(error);
        }
    };

    const handleCalculate = async () => {
        if (!logs.trim()) return;
        Vibration.vibrate(10);
        setLoading(true);
        const logLines = logs.split('\n').filter(line => line.trim().length > 0);
        const scoreResult = await calculateBoreholeScore(logLines);
        setResult(scoreResult);
        if (scoreResult.score) {
            Vibration.vibrate([0, 50, 50, 50]);
            await Database.saveScore(scoreResult.score, scoreResult.features || []);
        }
        setLoading(false);
    };

    const requestSmsPermission = async () => {
        try {
            const granted = await PermissionsAndroid.request(
                PermissionsAndroid.PERMISSIONS.READ_SMS,
                {
                    title: 'Borehole SMS Permission',
                    message: 'Borehole Engine needs access to your SMS to calculate your offline credit score.',
                    buttonNeutral: 'Ask Me Later',
                    buttonNegative: 'Cancel',
                    buttonPositive: 'OK',
                },
            );
            return granted === PermissionsAndroid.RESULTS.GRANTED;
        } catch (err) {
            console.warn(err);
            return false;
        }
    };

    const autoScanSms = async () => {
        const hasPermission = await requestSmsPermission();
        if (!hasPermission) {
            Alert.alert('Permission Denied', 'SMS permission is required for auto-scanning.');
            return;
        }

        setScanning(true);
        setLoading(true);

        const filter = {
            box: 'inbox',
            maxCount: 200,
        };

        SmsAndroid.list(
            JSON.stringify(filter),
            (fail: string) => {
                setScanning(false);
                setLoading(false);
                Alert.alert('Scan Failed', fail);
            },
            async (count: number, smsList: string) => {
                const messages = JSON.parse(smsList);
                const financialKeywords = [
                    'Confirmed', 'M-PESA', 'Airtel', 'HustlerFund',
                    'KCB', 'Equity', 'Okoa', 'Transaction ID'
                ];

                const filteredLogs = messages
                    .map((msg: any) => msg.body)
                    .filter((body: string) =>
                        financialKeywords.some(keyword => body.includes(keyword))
                    );

                if (filteredLogs.length === 0) {
                    setScanning(false);
                    setLoading(false);
                    Alert.alert('Scan Complete', 'No relevant financial logs detected in your inbox.');
                    return;
                }

                const scoreResult = await calculateBoreholeScore(filteredLogs);
                setResult(scoreResult);
                if (scoreResult.score) {
                    Vibration.vibrate([0, 50, 50, 50]);
                    await Database.saveScore(scoreResult.score, scoreResult.features || []);
                }
                setScanning(false);
                setLoading(false);
            },
        );
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
                    <TouchableOpacity onPress={loadHistory} style={styles.historyBtn}>
                        <Text style={styles.historyBtnText}>📜 History</Text>
                    </TouchableOpacity>
                </View>

                <View style={styles.card}>
                    <TouchableOpacity
                        style={styles.autoScanButton}
                        onPress={autoScanSms}
                        disabled={loading || scanning}
                    >
                        {scanning ? (
                            <ActivityIndicator color="#fff" />
                        ) : (
                            <Text style={styles.buttonText}>✨ Auto-Scan My Financial Health</Text>
                        )}
                    </TouchableOpacity>

                    <Text style={[styles.cardTitle, { marginTop: 20 }]}>Quick Test Personas</Text>
                    <View style={styles.chipRow}>
                        <TouchableOpacity style={[styles.chip, { backgroundColor: '#10B981' }]} onPress={() => setLogs(MOCK_PERSONAS.high_income)}>
                            <Text style={styles.chipText}>🤑 Rich</Text>
                        </TouchableOpacity>
                        <TouchableOpacity style={[styles.chip, { backgroundColor: '#F59E0B' }]} onPress={() => setLogs(MOCK_PERSONAS.average)}>
                            <Text style={styles.chipText}>😐 Avg</Text>
                        </TouchableOpacity>
                        <TouchableOpacity style={[styles.chip, { backgroundColor: '#EF4444' }]} onPress={() => setLogs(MOCK_PERSONAS.risky)}>
                            <Text style={styles.chipText}>🎰 Risky</Text>
                        </TouchableOpacity>
                    </View>

                    <Text style={[styles.cardTitle, { marginTop: 20 }]}>Or Paste Manually</Text>
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

                        <TouchableOpacity style={styles.verifyBtn} onPress={handleVerify}>
                            <Text style={styles.verifyBtnText}>🔐 Prove This Score (QR)</Text>
                        </TouchableOpacity>

                        {result.error && (
                            <Text style={styles.errorText}>Error: {result.error}</Text>
                        )}
                    </View>
                )}

                <View style={styles.footer}>
                    <Text style={styles.privacyNote}>🛡️ Privacy First: Your financial logs never leave this device.</Text>
                    <Text style={styles.privacyDetail}>🛡️ Privacy: Your SMS messages are processed locally by the Go-Engine and never uploaded to any server.</Text>
                </View>

            </ScrollView>

            <Modal visible={showHistory} animationType="slide">
                <SafeAreaView style={[styles.container, isDarkMode ? styles.darkBg : styles.lightBg]}>
                    <View style={styles.header}>
                        <Text style={styles.title}>Audit Log</Text>
                        <TouchableOpacity onPress={() => setShowHistory(false)} style={styles.closeBtn}>
                            <Text style={styles.closeBtnText}>Close</Text>
                        </TouchableOpacity>
                    </View>

                    <FlatList
                        data={history}
                        keyExtractor={(item) => item.id.toString()}
                        renderItem={({ item }) => (
                            <View style={styles.historyItem}>
                                <Text style={styles.historyDate}>
                                    {new Date(item.timestamp).toLocaleString()}
                                </Text>
                                <View style={styles.historyRow}>
                                    <Text style={styles.historyScore}>
                                        Score: {((item.score || 0) * 1000).toFixed(0)}
                                    </Text>
                                    {/* Parse features to show Income if possible */}
                                    <Text style={styles.historyDetail}>
                                        Logs: Encrypted
                                    </Text>
                                </View>
                            </View>
                        )}
                        contentContainerStyle={{ padding: 20 }}
                    />

                    <TouchableOpacity onPress={nukeHistory} style={styles.nukeBtn}>
                        <Text style={styles.nukeBtnText}>☢️ NUKE DATA</Text>
                    </TouchableOpacity>
                </SafeAreaView>
            </Modal>

            <Modal visible={showVerify} animationType="fade" transparent>
                <View style={styles.modalOverlay}>
                    <View style={styles.modalContent}>
                        <Text style={styles.modalTitle}>Verifiable Claim</Text>
                        <Text style={styles.modalSubtitle}>Scan to verify authenticity</Text>

                        <View style={styles.qrContainer}>
                            {cert?.signature ? (
                                <QRCode value={JSON.stringify(cert)} size={200} />
                            ) : (
                                <ActivityIndicator />
                            )}
                        </View>

                        <Text style={styles.certLabel}>Ed25519 Signature:</Text>
                        <Text style={styles.certHash} numberOfLines={2}>
                            {cert?.signature || 'Generating...'}
                        </Text>



                        <TouchableOpacity
                            style={styles.shareBtn}
                            onPress={handleShare}
                        >
                            <Text style={styles.shareBtnText}>📤 Share Proof</Text>
                        </TouchableOpacity>

                        <TouchableOpacity
                            style={styles.closeVerifyBtn}
                            onPress={() => setShowVerify(false)}
                        >
                            <Text style={styles.closeBtnText}>Done</Text>
                        </TouchableOpacity>
                    </View>
                </View>
            </Modal>
        </SafeAreaView >
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
    autoScanButton: {
        backgroundColor: '#10B981',
        padding: 18,
        borderRadius: 12,
        alignItems: 'center',
        elevation: 4,
        shadowColor: '#10B981',
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
        fontWeight: '600',
        marginBottom: 4,
    },
    privacyDetail: {
        fontSize: 10,
        color: '#94A3B8',
        fontWeight: '400',
        textAlign: 'center',
        paddingHorizontal: 40,
    },

    errorText: {
        color: '#EF4444',
        marginTop: 15,
        fontSize: 12,
        textAlign: 'center',
    },
    historyBtn: {
        marginTop: 10,
        backgroundColor: '#334155',
        paddingHorizontal: 12,
        paddingVertical: 6,
        borderRadius: 20,
    },
    historyBtnText: { color: '#fff', fontSize: 12, fontWeight: '700' },
    closeBtn: { position: 'absolute', right: 20, top: 35 },
    closeBtnText: { color: '#2196F3', fontSize: 16, fontWeight: '600' },
    historyItem: {
        backgroundColor: '#fff',
        padding: 15,
        borderRadius: 8,
        marginBottom: 10,
        elevation: 2,
    },
    historyDate: { fontSize: 12, color: '#94A3B8', marginBottom: 4 },
    historyRow: { flexDirection: 'row', justifyContent: 'space-between' },
    historyScore: { fontSize: 16, fontWeight: '700', color: '#0F172A' },
    historyDetail: { fontSize: 12, color: '#64748B' },
    nukeBtn: {
        backgroundColor: '#EF4444',
        margin: 20,
        padding: 15,
        borderRadius: 12,
        alignItems: 'center',
    },
    nukeBtnText: { color: '#fff', fontWeight: '800' },
    verifyBtn: {
        marginTop: 20,
        backgroundColor: '#334155',
        paddingVertical: 12,
        paddingHorizontal: 24,
        borderRadius: 24,
        width: '100%',
        alignItems: 'center',
    },
    verifyBtnText: { color: '#fff', fontWeight: '700', fontSize: 13 },
    modalOverlay: {
        flex: 1,
        backgroundColor: 'rgba(0,0,0,0.8)',
        justifyContent: 'center',
        alignItems: 'center',
        padding: 20,
    },
    modalContent: {
        width: '90%',
        backgroundColor: '#fff',
        borderRadius: 20,
        padding: 30,
        alignItems: 'center',
    },
    modalTitle: { fontSize: 24, fontWeight: '800', color: '#0F172A', marginBottom: 5 },
    modalSubtitle: { fontSize: 14, color: '#64748B', marginBottom: 20 },
    qrContainer: { padding: 20, backgroundColor: '#fff', borderRadius: 10, elevation: 5 },
    certLabel: { marginTop: 20, fontSize: 10, fontWeight: '700', color: '#94A3B8', textTransform: 'uppercase' },
    certHash: { fontSize: 10, color: '#334155', textAlign: 'center', marginTop: 5, fontFamily: 'monospace' },
    closeVerifyBtn: { marginTop: 10, padding: 10 },
    shareBtn: {
        marginTop: 20,
        backgroundColor: '#10B981',
        paddingVertical: 12,
        paddingHorizontal: 24,
        borderRadius: 24,
        width: '100%',
        alignItems: 'center',
    },
    shareBtnText: { color: '#fff', fontWeight: '700', fontSize: 13 },
    chipRow: { flexDirection: 'row', gap: 10, marginBottom: 10 },
    chip: { paddingVertical: 8, paddingHorizontal: 16, borderRadius: 20 },
    chipText: { color: '#fff', fontWeight: '700', fontSize: 12 },
});

export default DashboardScreen;
