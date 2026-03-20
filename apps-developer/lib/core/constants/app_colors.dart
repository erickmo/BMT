import 'package:flutter/material.dart';

class AppColors {
  AppColors._();

  // ─── Primary — Emerald Forest (Developer Edition) ────────────────────────────
  static const Color primary = Color(0xFF1A7A4A);
  static const Color primaryDark = Color(0xFF0C4A2B);
  static const Color primaryLight = Color(0xFF2EA878);

  // ─── Accent — Bright Emerald (replaces cyan) ─────────────────────────────────
  static const Color accent = Color(0xFF34D399);
  static const Color accentLight = Color(0xFF6EE7B7);

  // ─── Semantic ────────────────────────────────────────────────────────────────
  static const Color success = Color(0xFF16A34A);
  static const Color warning = Color(0xFFD97706);
  static const Color error = Color(0xFFDC2626);
  static const Color info = Color(0xFF0284C7);

  // ─── Background ──────────────────────────────────────────────────────────────
  static const Color background = Color(0xFFF0FDF4);
  static const Color surface = Color(0xFFFFFFFF);
  static const Color surfaceVariant = Color(0xFFDCFCE7);
  static const Color divider = Color(0xFFD1FAE5);

  // ─── Text ────────────────────────────────────────────────────────────────────
  static const Color textPrimary = Color(0xFF14532D);
  static const Color textSecondary = Color(0xFF4B7C5F);
  static const Color textHint = Color(0xFF86AB9A);
  static const Color textOnPrimary = Color(0xFFFFFFFF);

  // ─── Sidebar (Dark Forest) ───────────────────────────────────────────────────
  static const Color sidebarBg = Color(0xFF071A0D);
  static const Color sidebarActive = Color(0xFF1A7A4A);
  static const Color sidebarText = Color(0xFF6EE7B7);
  static const Color sidebarTextActive = Color(0xFFFFFFFF);

  // ─── Code (Dark Forest Theme) ────────────────────────────────────────────────
  static const Color codeBg = Color(0xFF071A0D);
  static const Color codeText = Color(0xFF6EE7B7);

  // ─── Gradients ───────────────────────────────────────────────────────────────
  static const LinearGradient headerGradient = LinearGradient(
    colors: [Color(0xFF071A0D), Color(0xFF1A7A4A)],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );
}
