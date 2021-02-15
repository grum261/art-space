from datetime import datetime
from typing import List, Optional, Tuple
from pydantic import BaseModel, EmailStr, AnyUrl, FilePath


class ModifiedAtMixin(BaseModel):
    created_at: Optional[datetime] = datetime.utcnow()
    updated_at: Optional[datetime] = datetime.utcnow()


class CompositionBase(BaseModel):
    composition: FilePath
    name: Optional[str] = 'Unnamed'
    genre: Optional[str] = None
    description: Optional[str] = None


class ArtCreate(CompositionBase):
    pass


class Art(CompositionBase):
    author_id: int


class AudioBase(CompositionBase):
    cover: Optional[FilePath] = None
    producer: Optional[str] = None
    featured: Optional[str] = None


class AudioCreate(AudioBase):
    pass


class Audio(AudioCreate):
    author_id: int


class ProfileBase(BaseModel):
    avatar: Optional[str] = None  # Optional[FilePath]
    activity: Optional[str] = None
    bio: Optional[str] = None
    social_networks: Optional[str] = None  # Optional[AnyUrl]


class ProfileCreate(ProfileBase):
    pass


class ProfileOut(ProfileBase):
    id: int


class UserBase(BaseModel):
    username: str
    email: EmailStr


class UserCreate(UserBase):
    password: str


class UserOut(UserBase, ModifiedAtMixin):
    id: int
    is_active: Optional[bool] = True
    profile: List[ProfileOut] = []
    # composition: Tuple[Art, Audio] = []

    class Config:
        orm_mode = True
