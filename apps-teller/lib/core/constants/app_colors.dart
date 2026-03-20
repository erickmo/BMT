import 'package:flutter/material.dart';

class AppColors {
  AppColors._();

  // ─── Primary — Emerald Pesantren ────────────────────────────────────────────
  static const Color primary = Color(0xFF1A7A4A);
  static const Color primaryDark = Color(0xFF0C4A2B);
  static const Color primaryMid = Color(0xFF2EA878);
  static const Color primaryLight = Color(0xFF4ADE80);
  static const Color primaryPale = Color(0xFFDCFCE7);

  // ─── Gold Accent (Islamic) ───────────────────────────────────────────────────
  static const Color gold = Color(0xFFC9A84C);
  static const Color goldLight = Color(0xFFF0CB6B);

  // ─── Secondary ───────────────────────────────────────────────────────────────
  static const Color secondary = Color(0xFF0284C7);
  static const Color secondaryLight = Color(0xFF38BDF8);

  // ─── Background ──────────────────────────────────────────────────────────────
  static const Color background = Color(0xFFF0FDF4);
  static const Color surface = Color(0xFFFFFFFF);
  static const Color cardBg = Color(0xFFFFFFFF);
  static const Color surfaceVariant = Color(0xFFDCFCE7);
  static const Color sidebarBg = Color(0xFF052E16);

  // ─── Text ────────────────────────────────────────────────────────────────────
  static const Color textPrimary = Color(0xFF14532D);
  static const Color textSecondary = Color(0xFF4B7C5F);
  static const Color textHint = Color(0xFF86AB9A);
  static const Color textOnPrimary = Color(0xFFFFFFFF);
  static const Color textOnSidebar = Color(0xFF6EE7B7);

  // ─── Status ──────────────────────────────────────────────────────────────────
  static const Color success = Color(0xFF16A34A);
  static const Color error = Color(0xFFDC2626);
  static const Color warning = Color(0xFFD97706);
  static const Color info = Color(0xFF0284C7);

  // ─── Transaksi ───────────────────────────────────────────────────────────────
  static const Color setor = Color(0xFF16A34A);
  static const Color tarik = Color(0xFFDC2626);

  // ─── Border ──────────────────────────────────────────────────────────────────
  static const Color border = Color(0xFFBBF7D0);
  static const Color divider = Color(0xFFD1FAE5);

  // ─── Gradients ───────────────────────────────────────────────────────────────
  static const LinearGradient headerGradient = LinearGradient(
    colors: [Color(0xFF0C4A2B), Color(0xFF1A7A4A)],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );

  static const LinearGradient cardGradient = LinearGradient(
    colors: [Color(0xFF1A7A4A), Color(0xFF2EA878)],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );
}
