from datetime import datetime
from typing import List, Optional
from pydantic import BaseModel, EmailStr, AnyUrl, FilePath


class ModifiedAtMixin(BaseModel):
    created_at: datetime
    updated_at: datetime


class CompositionBase(BaseModel):
    composition: [FilePath]
    name: Optional[str] = 'Unnamed'
    genre: Optional[str] = None
    description: Optional[str] = None


class ArtCreate(CompositionBase):
    pass


class Art(CompositionBase):
    author_id: int

    class Config:
        orm_mode = True


class AudioCreate(CompositionBase):
    cover: Optional[FilePath] = None
    producer: Optional[str] = None
    featured: Optional[str] = None


class Audio(AudioCreate):
    author_id: int

    class Config:
        orm_mode = True


class Composition(CompositionBase):
    art: Optional[CompositionBase]


class ProfileBase(BaseModel, ModifiedAtMixin):
    avatar: Optional[FilePath] = None
    activity: Optional[str] = None
    bio: Optional[str] = None
    social_networks: Optional[AnyUrl] = None


class ProfileCreate(ProfileBase):
    pass


class Profile(ProfileBase):
    id: int
    user_id: int

    class Config:
        orm_mode = True


class UserBase(BaseModel):
    username: str
    email: EmailStr


class UserCreate(UserBase):
    password: str


class User(UserBase, ModifiedAtMixin):
    id: int
    is_active: bool
    profile: List[Profile] = []
    composition: List[Art, Audio] = []

    class Config:
        orm_mode = True
