from datetime import datetime
from typing import List
from sqlalchemy import Boolean, Column, Integer, String, DateTime, Text, ForeignKey
from sqlalchemy.orm import relationship
from sqlalchemy.ext.declarative import declared_attr

from .db import Base


class IdMixin(object):
    @declared_attr
    def id(self):
        return Column(Integer, primary_key=True, index=True)


class ModifiedAtMixin(object):
    modified_at = Column(DateTime, default=datetime.utcnow())

    @declared_attr
    def created_at(self):
        return self.modified_at

    @declared_attr
    def update_at(self):
        return self.modified_at


class PathToFileMixin(object):
    @declared_attr
    def path_to_file(self):
        return Column(String(length=255))


class User(IdMixin, ModifiedAtMixin, Base):
    __table__ = 'users'

    username = Column(String, unique=True, index=True)
    email = Column(String, unique=True, index=True)
    hashed_password = Column(String)
    is_active = Column(Boolean, default=True)

    profile = relationship('Profile', back_populates='user', uselist=False)


class Profile(IdMixin, PathToFileMixin, Base):
    __table__ = 'profiles'

    activity = Column(String(length=255))
    bio = Column(Text(length=500))
    social_networks = Column(String(length=255))
    user_id = Column(Integer, ForeignKey('users.id'))

    user = relationship('User', back_populates='profile')
    composition: List = relationship('Composition', back_populates='author')


class Composition(IdMixin, ModifiedAtMixin, PathToFileMixin, Base):
    __abstract__ = True

    name = Column(String(length=255))
    genre = Column(String(length=255))
    description = Column(Text(length=500))
    author_id = Column(Integer, ForeignKey('users.id'))

    author = relationship('Profile', back_populates='composition')


class Art(Composition, Base):
    __table__ = 'arts'


class Audio(Composition, Base):
    __table__ = 'audios'

    cover = Column(String(length=255))
    producer = Column(String(length=255))
    featured = Column(String(length=255))
