package itswizard_m_berlinPreVersion

import (
	"encoding/json"
	"errors"
	itswizard_basic "github.com/itslearninggermany/itswizard_m_basic"
	imses "github.com/itslearninggermany/itswizard_m_imses"
	"log"
	"strings"
	"time"
)

type LusdUser struct {
	BenutzerDtGeaendertAm string `json:"benutzerDtGeaendertAm"`
	BenutzerGeburtsdatum  string `json:"benutzerGeburtsdatum"`
	BenutzerGlobalRolle   string `json:"benutzerGlobalRolle"`
	BenutzerNachname      string `json:"benutzerNachname"`

	BenutzerSchuleListe []BenutzerSchuleListe `json:"benutzerSchuleListe"`
	BenutzerVorname     string                `json:"benutzerVorname"`
	PersonUID           string                `json:"personUID"`
}

type BenutzerKlasseListe struct {
	BenutzerKlasseRolle string `json:"benutzerKlasseRolle"`
	KlasseUID           string `json:"klasseUID"`
}

type BenutzerKursListe struct {
	BenutzerKursRolle string `json:"benutzerKursRolle"`
	KursUID           string `json:"kursUID"`
}

type BenutzerSchuleListe struct {
	BenutzerKlasseListe []BenutzerKlasseListe `json:"benutzerKlasseListe"`
	BenutzerKursListe   []BenutzerKursListe   `json:"benutzerKursListe"`
	BenutzerSchuleRolle string                `json:"benutzerSchuleRolle"`
	BenutzerStatus      string                `json:"benutzerStatus"`
	SchuleUID           string                `json:"schuleUID"`
}

func personsParse(out []byte, err error) (lusdUsers []LusdUser, err2 error) {
	if err != nil {
		return nil, err
	}
	err2 = json.Unmarshal(out, &lusdUsers)
	return lusdUsers, err2
}

func (p *BlusdConnection) GetAllPersons() (lusdUsers []LusdUser, out []byte, err error) {
	out, err = p.callAPI(time.Now(), "", users, false)
	if err != nil {
		return
	}
	lusdUsers, err = personsParse(out, err)
	return
}

/*
Returns all new Persons created in the last 5 minutes
*/
func (p *BlusdConnection) GetAllPersonsCreatedScinse(minutes time.Duration) (lusdUsers []LusdUser, err error) {
	newT := time.Now().Add(-time.Minute * minutes)
	log.Println(newT.Format(timeLayout))
	return personsParse(p.callAPI(newT, "", users, true))
}

/*
Returns all Persons from a school
*/
func (p *BlusdConnection) GetAllPersonsFromSchool(sid string) (lusdUsers []LusdUser, err error) {
	return personsParse(p.callAPI(time.Now(), sid, users, false))
}

/*
Returns all new Persons from a school created in the last 5 minutes
*/
func (p *BlusdConnection) GetAllPersonsFromSchoolScince(sid string, minutes time.Duration) (lusdUsers []LusdUser, err error) {
	newT := time.Now().Add(-time.Minute * minutes)
	log.Println(newT.Format(timeLayout))
	return personsParse(p.callAPI(newT, sid, users, true))
}

/*
Returns all Persons from a school created since a given time
*/
func (p *BlusdConnection) GetAllPersonsFromSchoolScinceAGivenTime(sid string, t time.Time) (lusdUsers []LusdUser, err error) {
	return personsParse(p.callAPI(t, sid, users, true))
}

/*
Returns all Persons created since a given time
*/
func (p *BlusdConnection) GetAllPersonsScinceAGivenTime(t time.Time) (lusdUsers []LusdUser, err error) {
	return personsParse(p.callAPI(t, "", users, true))
}

/*
Get details from one given Person
*/
func (p *BlusdConnection) GetDetailOfPerson(personID string) (lusdUsers LusdUser, err error) {
	tmp, err := personsParse(p.callAPI(time.Now(), "", user+personID, false))
	if err != nil {
		if len(tmp) < 1 {
			erro := err.Error() + "The person with the personID " + personID + " does not exist!"
			return lusdUsers, errors.New(erro)
		}
		return lusdUsers, err
	}
	if len(tmp) < 1 {
		return lusdUsers, errors.New("The person with the personID " + personID + " does not exist!")
	}
	return tmp[0], err
}

/*
Imports all schools to itslearning. When the school exist it will be only updated.
*/
func (p *BlusdConnection) ImportOrUpdateAllPersonsToItslearning(itsl *imses.Request) (logContent []blusdLog, err error) {
	for _, person := range p.AllPersons {
		/*
			Todo: Schauen, ob es die Schule schon gibt.
		*/
		log.Println(person.BenutzerGlobalRolle)
		role := "Guest"
		if person.BenutzerGlobalRolle == "Schüler" {
			role = "Student"
		}
		if person.BenutzerGlobalRolle == "Lehrkraft" {
			role = "Staff"
		}

		resp, err := itsl.CreatePerson(itswizard_basic.DbPerson15{
			SyncPersonKey: person.PersonUID,
			FirstName:     person.BenutzerVorname,
			LastName:      person.BenutzerNachname,
			Profile:       role,
		})
		if err != nil {
			log.Println(resp)
		}

		for _, msh := range person.BenutzerSchuleListe {
			for _, membershipClass := range msh.BenutzerKlasseListe {
				membershipRole := "Guest"
				if strings.Contains(membershipClass.BenutzerKlasseRolle, "Schüler") {
					membershipRole = "Learner"
				}
				if strings.Contains(membershipClass.BenutzerKlasseRolle, "schüler") {
					membershipRole = "Learner"
				}
				if strings.Contains(membershipClass.BenutzerKlasseRolle, "lehrer") {
					membershipRole = "Instructor"
				}
				if strings.Contains(membershipClass.BenutzerKlasseRolle, "Lehrer") {
					membershipRole = "Instructor"
				}
				resp, err := itsl.CreateMembership(membershipClass.KlasseUID, person.PersonUID, membershipRole)
				//todo: log
				if err != nil {
					log.Println(resp)
				}
				log.Println(membershipClass.KlasseUID)
			}
			for _, membershipCourse := range msh.BenutzerKursListe {
				membershipRole := "Guest"
				if strings.Contains(membershipCourse.BenutzerKursRolle, "Schüler") {
					membershipRole = "Learner"
				}
				if strings.Contains(membershipCourse.BenutzerKursRolle, "lchüler") {
					membershipRole = "Learner"
				}
				if strings.Contains(membershipCourse.BenutzerKursRolle, "lehrer") {
					membershipRole = "Instructor"
				}
				if strings.Contains(membershipCourse.BenutzerKursRolle, "Lehrer") {
					membershipRole = "Instructor"
				}
				resp, err := itsl.CreateMembership(membershipCourse.KursUID, person.PersonUID, membershipRole)
				//todo: log
				if err != nil {
					log.Println(resp)
				}
				log.Println(membershipCourse.KursUID)
			}
		}
	}
	return logContent, nil
}

func (p *LusdUser) GetAllSchools() (out string, err error) {
	m := []string{}
	for _, k := range p.BenutzerSchuleListe {
		m = append(m, k.SchuleUID)
	}
	b, err := json.Marshal(m)
	return string(b), err
	return out, err
}

func (p *LusdUser) GetAllClasses() (out string, err error) {
	m := make(map[string]string)

	for _, k := range p.BenutzerSchuleListe {
		for _, k2 := range k.BenutzerKlasseListe {
			m[k2.KlasseUID] = k2.BenutzerKlasseRolle
		}
	}
	b, err := json.Marshal(m)
	return string(b), err
}

func (p *LusdUser) GetAllCourses() (out string, err error) {
	m := make(map[string]string)

	for _, k := range p.BenutzerSchuleListe {
		for _, k2 := range k.BenutzerKursListe {
			m[k2.KursUID] = k2.BenutzerKursRolle
		}
	}
	b, err := json.Marshal(m)
	return string(b), err
}
