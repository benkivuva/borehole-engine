package com.mobileapp;

import com.facebook.react.bridge.ReactApplicationContext;
import com.facebook.react.bridge.ReactContextBaseJavaModule;
import com.facebook.react.bridge.ReactMethod;
import com.facebook.react.bridge.Promise;

// Import the Go-generated package (this will resolving after successful gomobile bind)
import mobile.Mobile;
import mobile.MobileEngine;

public class BoreholeModule extends ReactContextBaseJavaModule {
    private static MobileEngine engine;

    BoreholeModule(ReactApplicationContext context) {
        super(context);
        if (engine == null) {
            engine = Mobile.newMobileEngine();
        }
    }

    @Override
    public String getName() {
        return "BoreholeModule";
    }

    @ReactMethod
    public void calculateBoreholeScore(String jsonLogs, Promise promise) {
        try {
            String result = engine.calculateBoreholeScore(jsonLogs);

            promise.resolve(result);
        } catch (Exception e) {
            promise.reject("ERR_SCORE", e.getMessage());
        }
    }

    @ReactMethod
    public void generateSignedScore(double score, Promise promise) {
        try {
            String result = engine.generateSignedScore(score);
            promise.resolve(result);
        } catch (Exception e) {
            promise.reject("ERR_SIGN", e.getMessage());
        }
    }
}
