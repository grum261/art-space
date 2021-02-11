from django.contrib import admin

from .models import Art, Audio


@admin.register(Art)
class CompositionAdmin(admin.ModelAdmin):
    list_display = (
        'id',
        'created_at',
        'composition_author',
        'picture_name',
        'direction',
    )

@admin.register(Audio)
class AudioAdmin(admin.ModelAdmin):
    list_display = (
        'id',
        'created_at',
        'composition_author',
        'track_name',
        'genre',
        'producer',
        'featured',
    )
