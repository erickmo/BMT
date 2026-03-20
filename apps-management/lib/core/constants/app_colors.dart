import 'package:flutter/material.dart';

class AppColors {
  AppColors._();

  // ─── Primary — Emerald Pesantren ────────────────────────────────────────────
  static const Color primary = Color(0xFF1A7A4A);
  static const Color primaryDark = Color(0xFF0C4A2B);
  static const Color primaryLight = Color(0xFF2EA878);
  static const Color primaryPale = Color(0xFFDCFCE7);

  // ─── Secondary ───────────────────────────────────────────────────────────────
  static const Color secondary = Color(0xFF0284C7);
  static const Color secondaryDark = Color(0xFF0369A1);
  static const Color secondaryLight = Color(0xFF38BDF8);

  // ─── Gold Accent ─────────────────────────────────────────────────────────────
  static const Color gold = Color(0xFFC9A84C);

  // ─── Semantic ────────────────────────────────────────────────────────────────
  static const Color success = Color(0xFF16A34A);
  static const Color warning = Color(0xFFD97706);
  static const Color error = Color(0xFFDC2626);
  static const Color info = Color(0xFF0284C7);

  // ─── Neutral ─────────────────────────────────────────────────────────────────
  static const Color background = Color(0xFFF0FDF4);
  static const Color surface = Color(0xFFFFFFFF);
  static const Color surfaceVariant = Color(0xFFDCFCE7);
  static const Color divider = Color(0xFFD1FAE5);

  // ─── Text ────────────────────────────────────────────────────────────────────
  static const Color textPrimary = Color(0xFF14532D);
  static const Color textSecondary = Color(0xFF4B7C5F);
  static const Color textHint = Color(0xFF86AB9A);
  static const Color textOnPrimary = Color(0xFFFFFFFF);

  // ─── Sidebar ─────────────────────────────────────────────────────────────────
  static const Color sidebarBg = Color(0xFF052E16);
  static const Color sidebarActive = Color(0xFF1A7A4A);
  static const Color sidebarText = Color(0xFF6EE7B7);
  static const Color sidebarTextActive = Color(0xFFFFFFFF);

  // ─── Status Badge ────────────────────────────────────────────────────────────
  static const Color statusAktif = Color(0xFF16A34A);
  static const Color statusPending = Color(0xFFD97706);
  static const Color statusDitolak = Color(0xFFDC2626);
  static const Color statusBlokir = Color(0xFF6B7280);

  // ─── Gradients ───────────────────────────────────────────────────────────────
  static const LinearGradient headerGradient = LinearGradient(
    colors: [Color(0xFF0C4A2B), Color(0xFF1A7A4A)],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );
}
