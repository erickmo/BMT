import 'package:shared_preferences/shared_preferences.dart';

class LocalStorage {
  final SharedPreferences _prefs;

  static const String _keyIsLoggedIn = 'teller_is_logged_in';
  static const String _keyNamaTeller = 'teller_nama';
  static const String _keyRole = 'teller_role';
  static const String _keyCabangId = 'teller_cabang_id';

  LocalStorage(this._prefs);

  Future<void> setLoggedIn(bool v) => _prefs.setBool(_keyIsLoggedIn, v);
  bool get isLoggedIn => _prefs.getBool(_keyIsLoggedIn) ?? false;

  Future<void> setNamaTeller(String nama) =>
      _prefs.setString(_keyNamaTeller, nama);
  String? get namaTeller => _prefs.getString(_keyNamaTeller);

  Future<void> setRole(String role) => _prefs.setString(_keyRole, role);
  String? get role => _prefs.getString(_keyRole);

  Future<void> setCabangId(String id) => _prefs.setString(_keyCabangId, id);
  String? get cabangId => _prefs.getString(_keyCabangId);

  Future<void> clear() => _prefs.clear();
}
