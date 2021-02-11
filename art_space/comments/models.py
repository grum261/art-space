from django.db import models


class AbstractComment(models.Model):
    '''Абстрактный комментарий'''
    
    comment = models.CharField(max_length=255)
    comment_author = models.ForeignKey(to='user.Profile', on_delete=models.DO_NOTHING, default='self')

    class Meta:
        abstract = True

class ArtComment(AbstractComment):
    '''Комментарий под картиной пользователя'''

    art = models.ForeignKey(to='composition.Art', on_delete=models.CASCADE)

    class Meta:
        db_table = 'comment_on_art'

class AudioComment(AbstractComment):
    '''Комментарий под треком пользователя'''

    audio = models.ForeignKey(to='composition.Audio', on_delete=models.CASCADE)

    class Meta:
        db_table = 'comment_on_music'