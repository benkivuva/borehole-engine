import SQLite from 'react-native-sqlite-storage';
import * as Keychain from 'react-native-keychain';

SQLite.enablePromise(true);

const DB_NAME = 'borehole_secure.db';
const KEYCHAIN_SERVICE = 'borehole_db_key';

export interface AuditLog {
    id: number;
    timestamp: number;
    score: number;
    features: string; // JSON
}

class DatabaseService {
    private db: SQLite.SQLiteDatabase | null = null;

    /**
     * Initialize Connection with Encryption Key from Keystore
     */
    async init(): Promise<void> {
        if (this.db) return;

        try {
            const key = await this.getOrGenerateKey();

            this.db = await SQLite.openDatabase({
                name: DB_NAME,
                location: 'default',
                key: key,
            });

            await this.createTables();
            console.log('🔒 Encrypted Database Initialized');
        } catch (error) {
            console.error('Database Init Error:', error);
            throw error;
        }
    }

    /**
     * Secure Key Management
     */
    private async getOrGenerateKey(): Promise<string> {
        try {
            // Check KeyStore
            const credentials = await Keychain.getGenericPassword({ service: KEYCHAIN_SERVICE });

            if (credentials) {
                return credentials.password;
            }

            const newKey = this.generateRandomKey();

            await Keychain.setGenericPassword('db_user', newKey, { service: KEYCHAIN_SERVICE });

            return newKey;
        } catch (err) {
            console.error('KeyStore Error:', err);
            throw new Error('Failed to access Concierge Keystore');
        }
    }

    private generateRandomKey(): string {
        const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()';
        let result = '';
        for (let i = 0; i < 64; i++) {
            result += chars.charAt(Math.floor(Math.random() * chars.length));
        }
        return result;
    }

    private async createTables() {
        if (!this.db) return;

        const schema = `
      CREATE TABLE IF NOT EXISTS audit_logs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        timestamp INTEGER NOT NULL,
        score REAL NOT NULL,
        features TEXT NOT NULL
      );
    `;

        await this.db.executeSql(schema);
    }

    /**
     * Public API
     */

    async saveScore(score: number, features: number[]): Promise<void> {
        if (!this.db) await this.init();

        const jsonFeatures = JSON.stringify(features);
        const timestamp = Date.now();

        await this.db!.executeSql(
            `INSERT INTO audit_logs (timestamp, score, features) VALUES (?, ?, ?)`,
            [timestamp, score, jsonFeatures]
        );
    }

    async getHistory(limit: number = 30): Promise<AuditLog[]> {
        if (!this.db) await this.init();

        const results = await this.db!.executeSql(
            `SELECT * FROM audit_logs ORDER BY timestamp DESC LIMIT ?`,
            [limit]
        );

        const logs: AuditLog[] = [];
        resultLoop: for (let i = 0; i < results[0].rows.length; i++) {
            logs.push(results[0].rows.item(i));
        }

        return logs;
    }

    async nukeData(): Promise<void> {
        if (!this.db) return;
        await this.db.executeSql('DELETE FROM audit_logs');
        await this.db.executeSql('VACUUM');
        console.log('☢️ Database Nuked');
    }
}

export const Database = new DatabaseService();
