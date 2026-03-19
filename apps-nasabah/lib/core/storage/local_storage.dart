import 'package:shared_preferences/shared_preferences.dart';

class LocalStorage {
  final SharedPreferences _prefs;

  static const String _keyIsLoggedIn = 'is_logged_in';
  static const String _keyNomorNasabah = 'nomor_nasabah';
  static const String _keyNamaNasabah = 'nama_nasabah';
  static const String _keyOnboardingDone = 'onboarding_done';

  LocalStorage(this._prefs);

  static Future<LocalStorage> create() async {
    final prefs = await SharedPreferences.getInstance();
    return LocalStorage(prefs);
  }

  Future<void> setLoggedIn(bool value) async {
    await _prefs.setBool(_keyIsLoggedIn, value);
  }

  bool get isLoggedIn => _prefs.getBool(_keyIsLoggedIn) ?? false;

  Future<void> setNomorNasabah(String nomor) async {
    await _prefs.setString(_keyNomorNasabah, nomor);
  }

  String? get nomorNasabah => _prefs.getString(_keyNomorNasabah);

  Future<void> setNamaNasabah(String nama) async {
    await _prefs.setString(_keyNamaNasabah, nama);
  }

  String? get namaNasabah => _prefs.getString(_keyNamaNasabah);

  Future<void> setOnboardingDone(bool value) async {
    await _prefs.setBool(_keyOnboardingDone, value);
  }

  bool get isOnboardingDone => _prefs.getBool(_keyOnboardingDone) ?? false;

  Future<void> clear() async {
    await _prefs.clear();
  }
}
