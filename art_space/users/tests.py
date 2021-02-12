from django.test import TestCase
from django.contrib.auth.models import User

from .models import Profile


class ProfileModelTest(TestCase):
    @classmethod
    def setUpTestData(cls):
        User.objects.create(email='a@a.com', username='a', password='qwertyasdyui')

    def test_user_profile_relation(self):
        print(User.objects.last().profile)
