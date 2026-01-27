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
$env:ANDROID_HOME = "C:\Users\ADMIN\AppData\Local\Android\sdk"; $env:ANDROID_NDK_HOME = "C:\Users\ADMIN\AppData\Local\Android\sdk\ndk\26.1.10909125"; $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User"); gomobile bind -target=android -o MobileApp/android/app/libs/borehole.aar ./pkg/mobile
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
1.  **Paste Logs**: Copy and paste raw M-Pesa, Airtel, or Bank SMS logs into the input field.
2.  **Calculate**: Tap **"Calculate Edge Score"**.
3.  **View Result**: The app calls the Go-Mobile engine via the JNI bridge and displays your **Borehole Index** (0-1000) and feature breakdown instantly, entirely offline.
