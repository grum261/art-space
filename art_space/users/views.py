from rest_framework import generics
from rest_framework.permissions import AllowAny, IsAuthenticated

from django.contrib.auth.models import User

from .serializers import RegisterUserSerializer, ProfileSerializer, ChangeUserPasswordSerializer
from .models import Profile


class RegisterUserView(generics.CreateAPIView):
    permission_classes = [AllowAny, ]
    queryset = User.objects.all()
    serializer_class = RegisterUserSerializer


class ChangeUserPasswordView(generics.UpdateAPIView):
    permission_classes = [IsAuthenticated, ]
    queryset = User.objects.all()
    serializer_class = ChangeUserPasswordSerializer


class ListProfileView(generics.ListAPIView):
    permission_classes = [AllowAny, ]
    queryset = Profile.objects.all()
    serializer_class = ProfileSerializer


class CreateProfileView(generics.CreateAPIView):
    permission_classes = [IsAuthenticated, ]
    queryset = Profile.objects.all()
    serializer = ProfileSerializer


class RetrieveUpdateDestroyProfileView(generics.RetrieveUpdateDestroyAPIView):
    permission_classes = [IsAuthenticated, ]
    queryset = Profile.objects.all()
    serializer = ProfileSerializer