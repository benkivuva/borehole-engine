import { NativeModules } from 'react-native';

const { BoreholeModule } = NativeModules;

export interface ScoreResult {
  Score: number;
  Features: number[];
  TxnCount: number;
  error?: string;
}

export const calculateBoreholeScore = async (logs: string[]): Promise<ScoreResult> => {
  try {
    const jsonLogs = JSON.stringify(logs);
    const resultJson = await BoreholeModule.calculateScore(jsonLogs);
    return JSON.parse(resultJson);
  } catch (error) {
    console.error('Borehole Scoring Error:', error);
    return {
      Score: 0,
      Features: [],
      TxnCount: 0,
      error: error instanceof Error ? error.message : 'Unknown error',
    };
  }
};

export default BoreholeModule;
