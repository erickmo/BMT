// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'login_response_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

LoginResponseModel _$LoginResponseModelFromJson(Map<String, dynamic> json) =>
    LoginResponseModel(
      id: json['id'] as String,
      nama: json['nama'] as String,
      username: json['username'] as String,
      role: json['role'] as String,
      bmtId: json['bmt_id'] as String,
      cabangId: json['cabang_id'] as String?,
      accessToken: json['access_token'] as String,
      refreshToken: json['refresh_token'] as String,
    );

Map<String, dynamic> _$LoginResponseModelToJson(LoginResponseModel instance) =>
    <String, dynamic>{
      'id': instance.id,
      'nama': instance.nama,
      'username': instance.username,
      'role': instance.role,
      'bmt_id': instance.bmtId,
      'cabang_id': instance.cabangId,
      'access_token': instance.accessToken,
      'refresh_token': instance.refreshToken,
    };
