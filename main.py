from typing import List
from fastapi import FastAPI

from asyncpg.exceptions import UniqueViolationError
from database.db import database
from database.schemas import UserOut, UserCreate, ProfileOut, ProfileCreate
from database.models import User, Profile

from sqlalchemy.orm import Session


app = FastAPI()


@app.on_event('startup')
async def startup():
    await database.connect()


@app.on_event('shutdown')
async def shutdown():
    await database.disconnect()


@app.get('/users/', response_model=List[UserOut])
async def read_users():
    query = User.__table__.select()

    return await database.fetch_all(query)


@app.post('/users/new/', response_model=UserOut)
async def create_user(user: UserCreate):
    query = User.__table__.insert().values(username=user.username,
                                           email=user.email, password=user.password)
    # Пока неясно как отловить говно, чтобы оно айдишник не писало, если в базе уже есть юзер
    try:
        last_record_id = await database.execute(query)
    except UniqueViolationError:
        raise UniqueViolationError
    return {**user.dict(), 'id': last_record_id}


@app.get('/profiles/', response_model=List[ProfileOut])
async def read_profiles():
    query = Profile.__table__.select()
    return await database.fetch_all(query)


@app.post('/profiles/new/', response_model=ProfileOut)
async def create_profile(profile: ProfileCreate):
    query = Profile.__table__.insert().values(avatar=profile.avatar, activity=profile.activity,
                                              bio=profile.bio, social_networks=profile.social_networks)
    last_record_id = await database.execute(query)

    return {**profile.dict(), 'id': last_record_id}
