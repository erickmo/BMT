class Validators {
  Validators._();

  static String? required(String? value, [String fieldName = 'Field']) {
    if (value == null || value.trim().isEmpty) {
      return '$fieldName tidak boleh kosong';
    }
    return null;
  }

  static String? password(String? value) {
    if (value == null || value.isEmpty) return 'Password tidak boleh kosong';
    return null;
  }

  static String? email(String? value) {
    if (value == null || value.isEmpty) return 'Email tidak boleh kosong';
    final emailRegex = RegExp(r'^[^@]+@[^@]+\.[^@]+$');
    if (!emailRegex.hasMatch(value)) return 'Format email tidak valid';
    return null;
  }

  static String? phone(String? value) {
    if (value == null || value.isEmpty) return 'Nomor telepon tidak boleh kosong';
    final phoneRegex = RegExp(r'^[0-9+\-\s]{8,15}$');
    if (!phoneRegex.hasMatch(value)) return 'Format nomor telepon tidak valid';
    return null;
  }

  static String? nominal(String? value) {
    if (value == null || value.isEmpty) return 'Nominal tidak boleh kosong';
    final cleaned = value.replaceAll(RegExp(r'[^\d]'), '');
    final amount = int.tryParse(cleaned);
    if (amount == null || amount <= 0) return 'Nominal harus lebih dari 0';
    return null;
  }

  static String? nik(String? value) {
    if (value == null || value.isEmpty) return 'NIK tidak boleh kosong';
    if (value.length != 16) return 'NIK harus 16 digit';
    if (!RegExp(r'^\d{16}$').hasMatch(value)) return 'NIK hanya boleh angka';
    return null;
  }
}
