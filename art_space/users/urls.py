from django.urls import path

from .views import (
    RegisterUserView, ListProfileView, 
    RetrieveUpdateDestroyProfileView, CreateProfileView, 
    ChangeUserPasswordView,
)


urlpatterns = [
    path('users/register/', RegisterUserView.as_view()),
    path('users/change_password/<int:pk>', ChangeUserPasswordView.as_view()),
    path('users/profiles/<int:pk>/', RetrieveUpdateDestroyProfileView.as_view()),
    path('users/profiles/', ListProfileView.as_view()),
    path('users/profiles/create/', CreateProfileView.as_view()),
]
