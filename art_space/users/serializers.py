from rest_framework import serializers

from django.contrib.auth.models import User
from django.contrib.auth.password_validation import validate_password
from django.core import exceptions

from .models import Profile


class RegisterUserSerializer(serializers.ModelSerializer):
    '''Сериалайзер пользователя'''

    class Meta:
        model = User
        fields = ('username', 'email', 'password')
        extra_kwargs = {
            'email': {'required': True},
            'password': {
                'write_only': True, 
                'required': True,
                'style': {'input_type': 'password'}
            },
        }

    def validate(self, attrs):
        user = User(**attrs)
        password = attrs.get('password')

        try:
            validate_password(password, user)
        except exceptions.ValidationError as error:
            serializer_error = serializers.as_serializer_error(error)
            
            raise serializers.ValidationError(
                {'password': serializer_error['non_field_errors']}
            )

        return attrs

    def create(self, validated_data):
        return User.objects.create_user(**validated_data)


class ChangeUserPasswordSerializer(serializers.ModelSerializer):
    '''Сериалайзер изменения'''

    old_password = serializers.CharField(write_only=True, required=True, style={'input': 'password'})
    new_password = serializers.CharField(
        write_only=True, required=True, 
        validators=[validate_password],
        style={'input': 'password'}
    )

    class Meta:
        model = User
        fields = ('old_password', 'new_password')

    def _validate_old_password(self, password):
        user = self.context['request'].user
        if user.check_password(password):
            return password
        raise serializers.ValidationError({'old_password': 'Пароль введен неправильно.'})

    def update(self, instance, validated_data):
        instance.set_password(validated_data['new_password'])
        instance.save()

        return instance


class ProfileSerializer(serializers.ModelSerializer):
    '''Сериалайзер пользовательского профиля'''

    class Meta:
        model = Profile
        fields = '__all__'