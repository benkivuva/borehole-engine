import React from 'react';
import {
    View,
    Text,
    TouchableOpacity,
    StyleSheet,
    useColorScheme,
    StatusBar,
    Dimensions,
} from 'react-native';
import Svg, { Path, Circle, Defs, LinearGradient, Stop } from 'react-native-svg';

interface WelcomeScreenProps {
    onStart: () => void;
}

const { width } = Dimensions.get('window');

const ShieldIcon = () => (
    <Svg width="120" height="140" viewBox="0 0 24 24" fill="none">
        <Defs>
            <LinearGradient id="grad" x1="0" y1="0" x2="0" y2="24">
                <Stop offset="0" stopColor="#3B82F6" stopOpacity="1" />
                <Stop offset="1" stopColor="#1E40AF" stopOpacity="1" />
            </LinearGradient>
        </Defs>
        <Path
            d="M12 2L3 7V13C3 18.5228 12 22 12 22C12 22 21 18.5228 21 13V7L12 2Z"
            fill="url(#grad)"
            stroke="#60A5FA"
            strokeWidth="1.5"
        />
        <Path
            d="M12 6L12 18"
            stroke="white"
            strokeWidth="2"
            strokeLinecap="round"
            opacity="0.5"
        />
        <Path
            d="M8 12L12 16L16 8"
            stroke="white"
            strokeWidth="2.5"
            strokeLinecap="round"
            strokeLinejoin="round"
        />
    </Svg>
);

const WelcomeScreen: React.FC<WelcomeScreenProps> = ({ onStart }) => {
    const isDarkMode = useColorScheme() === 'dark';
    const bgStyle = isDarkMode ? styles.darkBg : styles.lightBg;
    const textStyle = isDarkMode ? styles.darkText : styles.lightText;

    return (
        <View style={[styles.container, bgStyle]}>
            <StatusBar barStyle={isDarkMode ? 'light-content' : 'dark-content'} />

            <View style={styles.content}>
                <View style={styles.iconContainer}>
                    <ShieldIcon />
                    <View style={styles.glow} />
                </View>

                <Text style={[styles.title, textStyle]}>Borehole</Text>
                <Text style={styles.subtitle}>Sovereign Credit Infrastructure</Text>

                <View style={styles.featureContainer}>
                    <FeatureItem text="📱 100% Offline Scoring" />
                    <FeatureItem text="🔒 Private & Encrypted Vault" />
                    <FeatureItem text="🛡️ Cryptographically Verified" />
                </View>

                <TouchableOpacity
                    style={styles.button}
                    onPress={onStart}
                    activeOpacity={0.8}
                >
                    <Text style={styles.buttonText}>Launch Engine</Text>
                </TouchableOpacity>

                <Text style={styles.footer}>v1.0.0 • Powered by Go Mobile</Text>
            </View>
        </View>
    );
};

const FeatureItem = ({ text }: { text: string }) => (
    <View style={styles.featureRow}>
        <View style={styles.dot} />
        <Text style={styles.featureText}>{text}</Text>
    </View>
);

const styles = StyleSheet.create({
    container: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
    },
    content: {
        width: '100%',
        paddingHorizontal: 40,
        alignItems: 'center',
    },
    lightBg: { backgroundColor: '#F8FAFC' },
    darkBg: { backgroundColor: '#0F172A' },
    lightText: { color: '#0F172A' },
    darkText: { color: '#F1F5F9' },

    iconContainer: {
        marginBottom: 40,
        alignItems: 'center',
        justifyContent: 'center',
    },
    glow: {
        position: 'absolute',
        width: 60,
        height: 60,
        borderRadius: 30,
        backgroundColor: '#3B82F6',
        opacity: 0.15,
        transform: [{ scale: 2.5 }],
        zIndex: -1,
    },

    title: {
        fontSize: 42,
        fontWeight: '900',
        letterSpacing: -1,
        marginBottom: 8,
        fontVariant: ['small-caps'],
    },
    subtitle: {
        fontSize: 16,
        color: '#64748B',
        fontWeight: '500',
        letterSpacing: 0.5,
        textTransform: 'uppercase',
        marginBottom: 60,
    },

    featureContainer: {
        alignSelf: 'stretch',
        marginBottom: 60,
        paddingHorizontal: 10,
    },
    featureRow: {
        flexDirection: 'row',
        alignItems: 'center',
        marginBottom: 16,
    },
    dot: {
        width: 6,
        height: 6,
        borderRadius: 3,
        backgroundColor: '#3B82F6',
        marginRight: 12,
    },
    featureText: {
        fontSize: 15,
        color: '#94A3B8',
        fontWeight: '500',
    },

    button: {
        backgroundColor: '#2563EB',
        paddingVertical: 18,
        paddingHorizontal: 32,
        borderRadius: 16,
        width: '100%',
        alignItems: 'center',
        shadowColor: '#2563EB',
        shadowOffset: { width: 0, height: 8 },
        shadowOpacity: 0.3,
        shadowRadius: 16,
        elevation: 8,
    },
    buttonText: {
        color: '#FFFFFF',
        fontSize: 17,
        fontWeight: '700',
        letterSpacing: 0.5,
    },

    footer: {
        marginTop: 40,
        fontSize: 11,
        color: '#CBD5E1',
        fontWeight: '500',
    },
});

export default WelcomeScreen;
