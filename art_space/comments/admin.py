from django.contrib import admin

from .models import ArtComment, AudioComment


@admin.register(ArtComment)
class ArtCommentAdmin(admin.ModelAdmin):
    list_display = (
        'id', 'comment_author', 'art'
    )


@admin.register(AudioComment)
class AudioCommentAdmin(admin.ModelAdmin):
    list_display = (
        'id', 'comment_author', 'audio',
    )