class Validators {
  Validators._();

  static String? required(String? value, {String? fieldName}) {
    if (value == null || value.trim().isEmpty) {
      return '${fieldName ?? 'Field ini'} wajib diisi';
    }
    return null;
  }

  static String? pin(String? value) {
    if (value == null || value.isEmpty) {
      return 'PIN wajib diisi';
    }
    if (value.length != 6) {
      return 'PIN harus 6 digit';
    }
    if (!RegExp(r'^\d{6}$').hasMatch(value)) {
      return 'PIN hanya boleh berisi angka';
    }
    return null;
  }

  static String? nomorNasabah(String? value) {
    if (value == null || value.trim().isEmpty) {
      return 'Nomor nasabah wajib diisi';
    }
    if (value.trim().length < 4) {
      return 'Nomor nasabah tidak valid';
    }
    return null;
  }

  static String? nominal(String? value) {
    if (value == null || value.isEmpty) {
      return 'Nominal wajib diisi';
    }
    final cleaned = value.replaceAll(RegExp(r'[^0-9]'), '');
    final amount = int.tryParse(cleaned);
    if (amount == null || amount <= 0) {
      return 'Nominal harus lebih dari 0';
    }
    return null;
  }

  static String? email(String? value) {
    if (value == null || value.isEmpty) return null; // optional
    final regex = RegExp(r'^[\w\.-]+@[\w\.-]+\.\w{2,}$');
    if (!regex.hasMatch(value)) {
      return 'Email tidak valid';
    }
    return null;
  }

  static String? phone(String? value) {
    if (value == null || value.isEmpty) return null; // optional
    final cleaned = value.replaceAll(RegExp(r'[^0-9]'), '');
    if (cleaned.length < 9 || cleaned.length > 15) {
      return 'Nomor telepon tidak valid';
    }
    return null;
  }
}
