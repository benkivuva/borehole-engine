import { NativeModules } from 'react-native';

const { BoreholeModule } = NativeModules;

export interface ScoreResult {
  score: number;
  features: number[];
  txn_count: number;
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
      score: 0,
      features: [],
      txn_count: 0,
      error: error instanceof Error ? error.message : 'Unknown error',
    };
  }
};

export default BoreholeModule;
