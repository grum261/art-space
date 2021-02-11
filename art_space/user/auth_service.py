import re
from django.http import JsonResponse
from django.contrib.auth.models import User


def _exists(value) -> bool:
    if value:
        return True
    return JsonResponse({'succes': False, 'message': f'Передайте параметр {value}'})

def _is_valid_username(username):
    if re.search(r'[^a-zA-Z0-9\@\+\.\-\_]', username):
        return JsonResponse({
            'success': False, 
            'message': 'Имя пользователя должно состоять только из цифр, букв и символов @/./+/-/_'
        })
    return True

def _is_valid_password_length(password):
    if len(password) < 8:
        return JsonResponse({'success': False, 'message': 'Слишком короткий пароль'})
    return True

def _is_confirm_password_equal_to_password(confirm_password, password):
    if confirm_password == password:
        return True
    return JsonResponse({'succes': False, 'message': 'Пароли не совпадают'})

def validate_form(request):
    username, password, confirm_password = tuple(
        map(request.POST.get(), ('username', 'password', 'confirm_password'))
    )

    if (
        all(map(_exists, (username, password, confirm_password))) and 
        _is_valid_username(username) and _is_valid_password_length(password) and 
        _is_confirm_password_equal_to_password(confirm_password, password)):

        User.objects.create(username=username, password=password)

        return JsonResponse({'success': True, 'username': username})
    else:
        return JsonResponse({'success': False, 'message': 'Проверьте правильность введенных значений'})
