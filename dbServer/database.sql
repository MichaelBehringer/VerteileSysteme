CREATE TABLE IF NOT EXISTS Player (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    Username VARCHAR(255) NOT NULL,
    Skin VARCHAR(255),
    Passwort VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS Highscore (
    Highscore INT,
    Player_ID INT,
    FOREIGN KEY (Player_ID) REFERENCES Player(ID)
);

CREATE TABLE IF NOT EXISTS GameServer (
    ID VARCHAR(36) PRIMARY KEY,
    Servername VARCHAR(255) NOT NULL,
    Port INT NOT NULL,
    PlayerCounter INT NOT NULL,
    LastSeen TIMESTAMP NOT NULL
);

DELIMITER // 
CREATE PROCEDURE InsertUpdateGameServer(
	IN p_ID VARCHAR(36),
    IN p_Servername VARCHAR(255),
    IN p_Port INT,
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
        INSERT INTO GameServer (ID, Servername, Port, PlayerCounter, LastSeen)
        VALUES (p_ID, p_Servername, p_Port, p_PlayerCounter, NOW());
    END IF;
END; //  
DELIMITER ;

CREATE or REPLACE VIEW ActiveGameServer AS
SELECT ID, Servername, Port, PlayerCounter
FROM GameServer
WHERE TIMESTAMPDIFF(SECOND, LastSeen, NOW()) <= 15;

INSERT INTO Player (Username, Skin, Passwort) VALUES('Berry', 'superSkin', '123123');
INSERT INTO Player (Username, Skin, Passwort) VALUES('Dj', 'okSkin', 'Dj');
INSERT INTO Player (Username, Skin, Passwort) VALUES('MMMarco', 'keinSkin', 'abc');

INSERT INTO Highscore (Highscore, Player_ID) VALUES(69, 1);