from django.db import models

from user.models import Activities


class AbstractComposition(models.Model):
    '''Абстрактное произведение'''

    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    composition_author = models.ForeignKey(
        to='user.Profile', on_delete=models.DO_NOTHING, 
        default='self', verbose_name='Composition author'
    )
    category = models.CharField(max_length=255, choices=Activities.choices)
    description = models.TextField(max_length=500, blank=True)

    class Meta:
        abstract = True

class Audio(AbstractComposition):
    '''Песня, загружаемая пользователем'''

    track = models.FileField(upload_to='audio/', blank=True)
    cover = models.ImageField(upload_to='cover/', blank=True)
    track_name = models.CharField(max_length=255, default='Unnamed')
    genre = models.CharField(max_length=255, blank=True)
    producer = models.CharField(max_length=255, blank=True)
    featured = models.CharField(max_length=255, blank=True)

    def __str__(self):
        return self.track_name

    class Meta:
        db_table = 'audio'

class Art(AbstractComposition):
    '''Картина, загружаемая пользователем'''

    picture = models.ImageField(upload_to='art/', blank=True)
    picture_name = models.CharField(max_length=255, default='Unnamed')
    direction = models.CharField(max_length=255, blank=True)

    def __str__(self):
        return self.picture_name

    class Meta:
        db_table = 'art'