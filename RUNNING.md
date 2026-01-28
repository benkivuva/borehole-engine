# How to Run Borehole Mobile

Follow these steps to build and run the Borehole Edge-Scoring application on your Android device/emulator.

### Prerequisites
- **Node.js**: Installed (v18+)
- **JDK 17**: Installed and in your PATH (Verified)
- **Android SDK & NDK**: Installed via Android Studio
- **Go**: Installed (v1.22+)
- **gomobile**: Installed (`go install golang.org/x/mobile/cmd/gomobile@latest`)

---

### Step 1: Generate the Go Bridge (.aar)
The mobile app needs the Go engine compiled into an Android library. 
Run this command which sets the required environment variables for your system:

```powershell
# Set SDK/NDK paths (Adjust to your local system)
$env:ANDROID_HOME = "C:\Users\ADMIN\AppData\Local\Android\sdk"
$env:ANDROID_NDK_HOME = "$env:ANDROID_HOME\ndk\26.1.10909125"

# Build the bridge (Target API 21 for broad compatibility)
gomobile bind -v -target=android -androidapi 21 -o MobileApp/android/app/libs/borehole.aar ./pkg/mobile
```





### Step 2: Install Mobile Dependencies
Navigate to the mobile app folder and install the necessary React Native packages:

```powershell
cd MobileApp
npm install
```

### Step 3: Start the Metro Bundler
The bundler serves your JavaScript code to the application. Keep this terminal open:

```powershell
npx react-native start
```

### Step 4: Build and Run on Android
In a **new** terminal, run the application on your connected device or emulator:

```powershell
cd MobileApp
npx react-native run-android
```

---

### Operating the App
1.  **✨ Auto-Scan (Recommended)**: Tap the emerald green **"Auto-Scan My Financial Health"** button. The app will securely read your SMS inbox, filter for financial logs, and calculate your score instantly.
2.  **Manual Entry**: You can still paste raw logs into the text field and tap **"Calculate Edge Score"**.
3.  **Privacy**: No data is uploaded. All parsing and scoring occurs within the Go-Engine on your device.

---

### Testing with ADB (Emulator Only)
If you are using an emulator, you can inject test signals using these commands:

```powershell
# 1. Income (M-Pesa QKJ series)
adb emu sms send Sarah "QKJ3XPYC5T Confirmed. You have received Ksh15,000.00 from SARAH JANE on 25/1/26."

# 2. Debt (Okoa Jahazi Combined Signal)
adb emu sms send 444 "You have received Ksh 100.00 Okoa Jahazi. Your Okoa debt is Ksh 110.00."

# 3. Repayment (Hustler Fund "Sent" keyword)
adb emu sms send Hustler "Confirmed. You have sent Ksh2,000.00 to Hustler Fund on 20/1/26."
```

