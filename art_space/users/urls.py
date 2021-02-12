from django.urls import path

from .views import (
    RegisterUserView, ListProfileView, 
    RetrieveUpdateDestroyProfileView, CreateProfileView, 
    ChangeUserPasswordView,
)


urlpatterns = [
    path('users/registration/', RegisterUserView.as_view()),
    path('users/<int:pk>/change_password/', ChangeUserPasswordView.as_view()),
    path('users/profiles/<int:pk>/', RetrieveUpdateDestroyProfileView.as_view()),
    path('users/profiles/', ListProfileView.as_view()),
    path('users/profiles/new/', CreateProfileView.as_view()),
]
