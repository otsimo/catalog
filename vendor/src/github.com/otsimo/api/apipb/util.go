package apipb

import "gopkg.in/mgo.v2/bson"

//Child
func (c *Child) To24bytesId() {
	if len(c.Id) == 12 {
		c.Id = bson.ObjectId(c.Id.Hex())
	}
	if len(c.ParentId) == 12 {
		c.ParentId = bson.ObjectId(c.ParentId.Hex())
	}
	for _, ge := range c.Games {
		ge.To24bytesId()
	}
}

func (c *Child) To12bytesId() {
	if len(c.Id) == 24 {
		c.Id = bson.ObjectIdHex(string(c.Id))
	}
	if len(c.ParentId) == 24 {
		c.ParentId = bson.ObjectIdHex(string(c.ParentId))
	}
	for _, ge := range c.Games {
		ge.To12bytesId()
	}
}

//Profile
func (c *Profile) To24bytesId() {
	if len(c.Id) == 12 {
		c.Id = bson.ObjectId(c.Id.Hex())
	}
}

func (c *Profile) To12bytesId() {
	if len(c.Id) == 24 {
		c.Id = bson.ObjectIdHex(string(c.Id))
	}
}

//ChildGameEntry
func (c *ChildGameEntry) To24bytesId() {
	if len(c.Id) == 12 {
		c.Id = bson.ObjectId(c.Id.Hex())
	}
}

func (c *ChildGameEntry) To12bytesId() {
	if len(c.Id) == 24 {
		c.Id = bson.ObjectIdHex(string(c.Id))
	}
}

//Game
func (c *Game) To24bytesId() {
	if len(c.Id) == 12 {
		c.Id = bson.ObjectId(c.Id.Hex())
	}
	if len(c.OwnerId) == 12 {
		c.OwnerId = bson.ObjectId(c.OwnerId.Hex())
	}
}

func (c *Game) To12bytesId() {
	if len(c.Id) == 24 {
		c.Id = bson.ObjectIdHex(string(c.Id))
	}
	if len(c.OwnerId) == 24 {
		c.OwnerId = bson.ObjectIdHex(string(c.OwnerId))
	}
}

//GameRelease
func (c *GameRelease) To24bytesId() {
	if len(c.ReleaseId) == 12 {
		c.ReleaseId = bson.ObjectId(c.ReleaseId.Hex())
	}
	if len(c.GameId) == 12 {
		c.GameId = bson.ObjectId(c.GameId.Hex())
	}
	if len(c.ReleasedBy) == 12 {
		c.ReleasedBy = bson.ObjectId(c.ReleasedBy.Hex())
	}
	if len(c.ValidatedBy) == 12 {
		c.ValidatedBy = bson.ObjectId(c.ValidatedBy.Hex())
	}
}

func (c *GameRelease) To12bytesId() {
	if len(c.ReleaseId) == 24 {
		c.ReleaseId = bson.ObjectIdHex(string(c.ReleaseId))
	}
	if len(c.GameId) == 24 {
		c.GameId = bson.ObjectIdHex(string(c.GameId))
	}
	if len(c.ReleasedBy) == 24 {
		c.ReleasedBy = bson.ObjectIdHex(string(c.ReleasedBy))
	}
	if len(c.ValidatedBy) == 24 {
		c.ValidatedBy = bson.ObjectIdHex(string(c.ValidatedBy))
	}
}

//UploadToken
func (c *UploadToken) To24bytesId() {
	if len(c.GameId) == 12 {
		c.GameId = bson.ObjectId(c.GameId.Hex())
	}
	if len(c.UserId) == 12 {
		c.UserId = bson.ObjectId(c.UserId.Hex())
	}
}

func (c *UploadToken) To12bytesId() {
	if len(c.GameId) == 24 {
		c.GameId = bson.ObjectIdHex(string(c.GameId))
	}
	if len(c.UserId) == 24 {
		c.UserId = bson.ObjectIdHex(string(c.UserId))
	}
}
