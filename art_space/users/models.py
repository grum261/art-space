from django.db import models
from django.contrib.auth.models import User
from django.db.models.signals import post_save
from django.dispatch import receiver


class Activities(models.TextChoices):
    MUSIC = 'Music'
    ART = 'Art'

class Profile(models.Model):
    '''Пользовательский профиль'''

    user = models.OneToOneField(User, on_delete=models.CASCADE)
    location = models.CharField(max_length=50, blank=True)
    activity = models.CharField(max_length=255, choices=Activities.choices, blank=True)
    avatar = models.ImageField(upload_to='avatar/', null=True, blank=True)
    bio = models.TextField(max_length=500, blank=True)
    social_networks = models.URLField(max_length=255, blank=True)

    def __str__(self):
        return self.user.username

    class Meta:
        db_table = 'profile'


@receiver(post_save, sender=User)
def create_or_update_user_profile(sender, instance, created, **kwargs):
    if created:
        Profile.objects.create(user=instance)
    instance.profile.save()