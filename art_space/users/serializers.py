from rest_framework import serializers

from django.contrib.auth.models import User
from django.contrib.auth.password_validation import validate_password

from .models import Profile


class RegisterUserSerializer(serializers.ModelSerializer):
    '''Сериалайзер пользователя'''

    class Meta:
        model = User
        fields = ('username', 'email', 'password', 'date_joined')
        extra_kwargs = {
            #'email': {'required': True},
            'password': {'write_only': True}
        }

    def create(self, validated_data):
        return User.objects.create_user(username=validated_data['username'], password=validated_data['password'])


class ChangeUserPasswordSerializer(serializers.ModelSerializer):
    old_password = serializers.CharField(write_only=True, required=True)
    new_password = serializers.CharField(write_only=True, required=True, validators=[validate_password])

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