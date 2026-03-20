import 'package:flutter/material.dart';

class AppColors {
  AppColors._();

  // ─── Primary — Emerald Pesantren ────────────────────────────────────────────
  static const Color primary = Color(0xFF1A7A4A);
  static const Color primaryDark = Color(0xFF0C4A2B);
  static const Color primaryLight = Color(0xFF2EA878);

  // ─── Secondary ───────────────────────────────────────────────────────────────
  static const Color secondary = Color(0xFFC9A84C);
  static const Color success = Color(0xFF16A34A);
  static const Color warning = Color(0xFFD97706);
  static const Color error = Color(0xFFDC2626);

  // ─── Background ──────────────────────────────────────────────────────────────
  static const Color background = Color(0xFFF0FDF4);
  static const Color surface = Color(0xFFFFFFFF);
  static const Color surfaceVariant = Color(0xFFDCFCE7);
  static const Color divider = Color(0xFFD1FAE5);
  static const Color border = Color(0xFFBBF7D0);

  // ─── Text ────────────────────────────────────────────────────────────────────
  static const Color textPrimary = Color(0xFF14532D);
  static const Color textSecondary = Color(0xFF4B7C5F);
  static const Color textHint = Color(0xFF86AB9A);
  static const Color textOnPrimary = Color(0xFFFFFFFF);

  // ─── POS / Kasir ─────────────────────────────────────────────────────────────
  static const Color nfcActive = Color(0xFF4ADE80);
  static const Color nfcIdle = Color(0xFF6B7280);
  static const Color kasirBg = Color(0xFF052E16);

  // ─── Gradients ───────────────────────────────────────────────────────────────
  static const LinearGradient headerGradient = LinearGradient(
    colors: [Color(0xFF0C4A2B), Color(0xFF1A7A4A)],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );
}
