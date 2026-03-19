import 'package:json_annotation/json_annotation.dart';
import '../../domain/entities/user_entity.dart';

part 'login_response_model.g.dart';

@JsonSerializable()
class LoginResponseModel {
  final String id;
  final String nama;
  final String username;
  final String role;
  @JsonKey(name: 'bmt_id')
  final String bmtId;
  @JsonKey(name: 'cabang_id')
  final String? cabangId;
  @JsonKey(name: 'access_token')
  final String accessToken;
  @JsonKey(name: 'refresh_token')
  final String refreshToken;

  const LoginResponseModel({
    required this.id,
    required this.nama,
    required this.username,
    required this.role,
    required this.bmtId,
    this.cabangId,
    required this.accessToken,
    required this.refreshToken,
  });

  factory LoginResponseModel.fromJson(Map<String, dynamic> json) =>
      _$LoginResponseModelFromJson(json);

  Map<String, dynamic> toJson() => _$LoginResponseModelToJson(this);

  UserEntity toEntity() => UserEntity(
        id: id,
        nama: nama,
        username: username,
        role: role,
        bmtId: bmtId,
        cabangId: cabangId,
        accessToken: accessToken,
        refreshToken: refreshToken,
      );
}
