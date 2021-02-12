from rest_framework import generics
from rest_framework.permissions import AllowAny, IsAuthenticated

from django.contrib.auth.models import User

from .serializers import RegisterUserSerializer, ProfileSerializer, ChangeUserPasswordSerializer
from .models import Profile


class RegisterUserView(generics.CreateAPIView):
    '''Страница регистрации нового пользователя: /api/users/registration/'''
    permission_classes = [AllowAny, ]
    queryset = User.objects.all()
    serializer_class = RegisterUserSerializer


class ChangeUserPasswordView(generics.UpdateAPIView):
    '''
    Страница изменения пароля пользователя: /api/users/{user_id}/change_password/

    Создать и добавить права на изменения только самим собой
    '''
    permission_classes = [IsAuthenticated, ]
    queryset = User.objects.all()
    serializer_class = ChangeUserPasswordSerializer


class ListProfileView(generics.ListAPIView):
    '''Страница просмотра пользовательских профилей: /api/users/profiles/'''
    permission_classes = [AllowAny, ]
    queryset = Profile.objects.all()
    serializer_class = ProfileSerializer


class CreateProfileView(generics.CreateAPIView):
    '''Страница создания профиля зарегистрированного пользователя: /api/users/profiles/new/'''
    permission_classes = [IsAuthenticated, ]
    queryset = Profile.objects.all()
    serializer_class = ProfileSerializer


class RetrieveUpdateDestroyProfileView(generics.RetrieveUpdateDestroyAPIView):
    '''Страница изменения/удаления профиля пользователя: /api/users/profiles/{profile_id}/'''
    permission_classes = [IsAuthenticated, ]
    queryset = Profile.objects.all()
    serializer_class = ProfileSerializer