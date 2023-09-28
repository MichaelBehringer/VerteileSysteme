CREATE TABLE IF NOT EXISTS Player (
    ID VARCHAR(36) PRIMARY KEY,
    Username VARCHAR(255) NOT NULL,
    Gamename VARCHAR(255) NOT NULL,
    Skin VARCHAR(255) NOT NULL,
    Passwort VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS Highscore (
    Player_ID VARCHAR(36) NOT NULL,
    Server_ID VARCHAR(36) NOT NULL,
    Score INT NOT NULL,
    PRIMARY KEY (Player_ID, Server_ID)
);

CREATE TABLE IF NOT EXISTS GameServer (
    ID VARCHAR(36) PRIMARY KEY,
    Servername VARCHAR(255) NOT NULL,
    Servernumber INT NOT NULL,
    PlayerCounter INT NOT NULL,
    LastSeen TIMESTAMP NOT NULL
);

DELIMITER // 
CREATE PROCEDURE InsertUpdateGameServer(
	IN p_ID VARCHAR(36),
    IN p_Servername VARCHAR(255),
    IN p_Servernumber INT,
    IN p_PlayerCounter INT
)
BEGIN
    DECLARE v_RecordCount INT;

    -- Prüfen, ob die ID bereits vorhanden ist
    SELECT COUNT(*) INTO v_RecordCount FROM GameServer WHERE ID = p_ID;

    IF v_RecordCount > 0 THEN
        -- Eintrag aktualisieren, wenn die ID bereits vorhanden ist
        UPDATE GameServer
        SET
            PlayerCounter = p_PlayerCounter,
            LastSeen = NOW()
        WHERE ID = p_ID;
    ELSE
        -- Neuen Eintrag einfügen, wenn die ID nicht vorhanden ist
        INSERT INTO GameServer (ID, Servername, Servernumber, PlayerCounter, LastSeen)
        VALUES (p_ID, p_Servername, p_Servernumber, p_PlayerCounter, NOW());
    END IF;
END; //  
DELIMITER ;

CREATE or REPLACE VIEW ActiveGameServer AS
SELECT ID, Servername, Servernumber, PlayerCounter
FROM GameServer
WHERE TIMESTAMPDIFF(SECOND, LastSeen, NOW()) <= 15;

CREATE or REPLACE VIEW HighscoreList AS
SELECT p.Gamename as Username, max(h.Score) as Score
FROM Highscore h
INNER JOIN Player p on h.Player_ID = p.ID
GROUP BY p.ID
ORDER BY h.Score desc
LIMIT 5;

INSERT INTO Player (ID, Username, Gamename, Skin, Passwort) VALUES('00000000-0000-0000-0000-000000000001', 'Berry', 'Berry', 'red', 'abc123');
INSERT INTO Player (ID, Username, Gamename, Skin, Passwort) VALUES('00000000-0000-0000-0000-000000000002', 'Dj', 'Dj', 'green', 'smash');
INSERT INTO Player (ID, Username, Gamename, Skin, Passwort) VALUES('00000000-0000-0000-0000-000000000003', 'Marco', 'Marco', 'blue', 'kommentar');

INSERT INTO GameServer (ID, Servername, Servernumber, PlayerCounter, LastSeen) VALUES ('00000000-0000-0000-0000-000000000010','open-blowfish',1,0,'2023-09-28 10:45:39.000');
INSERT INTO GameServer (ID, Servername, Servernumber, PlayerCounter, LastSeen) VALUES ('00000000-0000-0000-0000-000000000020','summary-boa',2,0,'2023-09-28 12:06:27.000');
INSERT INTO GameServer (ID, Servername, Servernumber, PlayerCounter, LastSeen) VALUES ('00000000-0000-0000-0000-000000000030','green-dog',3,0,'2023-09-28 12:07:27.000');

INSERT INTO Highscore (Player_ID, Server_ID, Score) VALUES('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000010', 69);
INSERT INTO Highscore (Player_ID, Server_ID, Score) VALUES('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000020', 40);
INSERT INTO Highscore (Player_ID, Server_ID, Score) VALUES('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000020', 70);
INSERT INTO Highscore (Player_ID, Server_ID, Score) VALUES('00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000030', 2);
