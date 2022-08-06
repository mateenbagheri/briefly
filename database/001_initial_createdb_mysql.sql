CREATE TABLE Users (
	userID INT NOT NULL AUTO_INCREMENT,
	name VARCHAR(50),
	familyName VARCHAR(50),
	password VARCHAR(100),
	salt VARCHAR(100),
	email VARCHAR(100),
	CONSTRAINT PK_Users PRIMARY KEY(userID)
	CONSTRAINT UC_Users UNIQUE (email)
);

CREATE TABLE Collections (
	collectionID INT NOT NULL AUTO_INCREMENT,
	collectionName VARCHAR(50),
	userID INT,
	CONSTRAINT PK_Collections PRIMARY KEY(collectionID),
	CONSTRAINT FK_Collections_UserID FOREIGN KEY (userID)
		REFERENCES Users(userID)
		ON DELETE CASCADE
);

CREATE TABLE Links (
	linkID BIGINT NOT NULL AUTO_INCREMENT,
	link VARCHAR(100) NOT NULL,
	shortened VARCHAR(100) NOT NULL,
	expDate DATE NOT NULL,
	collectionID INT,
	hitNumbers BIGINT,
	CONSTRAINT PK_Links PRIMARY KEY(linkID),
	CONSTRAINT FK_Links_collectionID FOREIGN KEY (collectionID)
        REFERENCES Collections(collectionID)
        ON DELETE CASCADE
);
