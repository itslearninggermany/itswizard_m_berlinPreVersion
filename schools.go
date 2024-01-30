package itswizard_m_berlinPreVersion

import (
	"encoding/json"
	itswizard_basic "github.com/itslearninggermany/itswizard_m_basic"
	imses "github.com/itslearninggermany/itswizard_m_imses"
	"log"
	"time"
)

type LusdSchool struct {
	SchuleBezirkName     string   `json:"schuleBezirkName"`
	SchuleBezirkNummer   string   `json:"schuleBezirkNummer"`
	SchuleDtGeaendertAm  string   `json:"schuleDtGeaendertAm"`
	SchuleName           string   `json:"schuleName"`
	SchuleNummer         string   `json:"schuleNummer"`
	SchuleSchulformListe []string `json:"schuleSchulformListe"`
	SchuleTyp            string   `json:"schuleTyp"`
	SchuleUID            string   `json:"schuleUID"`
}

func schoolParse(out []byte, err error) (lusdSchools []LusdSchool, err2 error) {
	if err != nil {
		return lusdSchools, err
	}
	err2 = json.Unmarshal(out, &lusdSchools)
	return lusdSchools, err2
}

/*
Returns all Schools from lusd
*/
func (p *BlusdConnection) GetAllSchools() (lusdSchools []LusdSchool, out []byte, err error) {
	out, err = p.callAPI(time.Now(), "", school, false)
	if err != nil {
		return
	}
	lusdSchools, err = schoolParse(out, err)
	return
}

/*
Returns all new Courses created in the last 5 minutes
*/
func (p *BlusdConnection) GetAllSchoolsCreatedScinse(minutes time.Duration) (lusdSchools []LusdSchool, err error) {
	newT := time.Now().Add(-time.Minute * minutes)
	log.Println(newT.Format(timeLayout))
	return schoolParse(p.callAPI(newT, "", school, true))
}

/*
Returns all Courses from a school
*/
func (p *BlusdConnection) GetAllSchoolsFromSchool(sid string) (lusdSchools []LusdSchool, err error) {
	return schoolParse(p.callAPI(time.Now(), sid, school, false))
}

/*
Returns all new Courses from a school created in the last 5 minutes
*/
func (p *BlusdConnection) GetAllSchoolsFromSchoolScince(sid string, minutes time.Duration) (lusdSchools []LusdSchool, err error) {
	newT := time.Now().Add(-time.Minute * minutes)
	log.Println(newT.Format(timeLayout))
	return schoolParse(p.callAPI(newT, sid, school, true))
}

/*
Returns all Courses from a school created since a given time
*/
func (p *BlusdConnection) GetAllSchoolsFromSchoolScinceAGivenTime(sid string, t time.Time) (lusdSchools []LusdSchool, err error) {
	return schoolParse(p.callAPI(t, sid, school, true))
}

/*
Returns all Courses created since a given time
*/
func (p *BlusdConnection) GetAllSchoolsScinceAGivenTime(t time.Time) (lusdSchools []LusdSchool, err error) {
	return schoolParse(p.callAPI(t, "", school, true))
}

/*
Imports all schools to itslearning. When the school exist it will be only updated.
*/
func (p *BlusdConnection) ImportOrUpdateAllSchoolsToItslearning(itsl *imses.Request) (logContent []blusdLog, err error) {
	lusdSchools, _, err := p.GetAllSchools()
	if err != nil {
		return logContent, err
	}
	for _, school := range lusdSchools {
		/*
			Todo: Schauen, ob es die Schule schon gibt.
		*/
		tmp := school.SchuleNummer + " " + school.SchuleName
		res := ""
		var gro []string
		for i, r := range tmp {
			res = res + string(r)
			if i > 0 && (i+1)%64 == 0 {
				gro = append(gro, res)
				res = ""
			}
		}
		schoolname := res
		if len(gro) > 0 {
			schoolname = gro[0]
			log.Println("gro:", gro[0])
		}

		resp, err := itsl.CreateGroup(itswizard_basic.DbGroup15{
			SyncID:        school.SchuleUID,
			Name:          schoolname,
			ParentGroupID: "0",
		}, true)
		log.Println(schoolname)
		if err != nil {
			logContent = append(logContent, blusdLog{
				Success: true,
				Date:    time.Now(),
				Message: "School with id " + school.SchuleUID + " and name " + schoolname + " was not (!) created in itslearning. " + resp,
			})
		} else {
			logContent = append(logContent, blusdLog{
				Success: true,
				Date:    time.Now(),
				Message: "School with id " + school.SchuleUID + " and name " + schoolname + " was created in itslearning.",
			})
		}
		resp, err = itsl.CreateGroup(itswizard_basic.DbGroup15{
			SyncID:        school.SchuleUID + "-TEACHERS",
			Name:          "Lehrkräfte",
			ParentGroupID: school.SchuleUID,
		}, false)
		log.Println(school.SchuleUID + "-TEACHERS")

		if err != nil {
			logContent = append(logContent, blusdLog{
				Success: true,
				Date:    time.Now(),
				Message: "Group with id " + school.SchuleUID + "-TEACHERS " + "was not (!) created in itslearning. " + resp,
			})
		}
		resp, err = itsl.CreateGroup(itswizard_basic.DbGroup15{
			SyncID:        school.SchuleUID + "-STUDENTS",
			Name:          "Schüer",
			ParentGroupID: school.SchuleUID,
		}, false)
		log.Println(school.SchuleUID + "-STUDENTS")
		if err != nil {
			logContent = append(logContent, blusdLog{
				Success: true,
				Date:    time.Now(),
				Message: "Group with id " + school.SchuleUID + "-STUDENTS " + "was not (!) created in itslearning. " + resp,
			})
		}
		resp, err = itsl.CreateGroup(itswizard_basic.DbGroup15{
			SyncID:        school.SchuleUID + "-CLASSES",
			Name:          "Klassen",
			ParentGroupID: school.SchuleUID,
		}, false)
		log.Println(school.SchuleUID + "-CLASSES")
		if err != nil {
			logContent = append(logContent, blusdLog{
				Success: true,
				Date:    time.Now(),
				Message: "Group with id " + school.SchuleUID + "-CLASSES " + "was not (!) created in itslearning. " + resp,
			})
		}
		resp, err = itsl.CreateGroup(itswizard_basic.DbGroup15{
			SyncID:        school.SchuleUID + "-COURSES",
			Name:          "Kurse",
			ParentGroupID: school.SchuleUID,
		}, false)
		log.Println(school.SchuleUID + "-COURSES")
		if err != nil {
			logContent = append(logContent, blusdLog{
				Success: true,
				Date:    time.Now(),
				Message: "Group with id " + school.SchuleUID + "-COURSES " + "was not (!) created in itslearning. " + resp,
			})
		}
	}
	return logContent, nil
}

/*
TODO: Delete all classes which are not anymore in lusd.
*/
func (p *BlusdConnection) DeleteSchoolsInItslearning(itsl *imses.Request) (logContent []blusdLog, err error) {

	return logContent, nil
}

func (p *LusdSchool) SendToItslearning(itsl *imses.Request) (logContent []blusdLog, err error) {
	tmp := p.SchuleNummer + " " + p.SchuleName
	res := ""
	var gro []string
	for i, r := range tmp {
		res = res + string(r)
		if i > 0 && (i+1)%64 == 0 {
			gro = append(gro, res)
			res = ""
		}
	}
	schoolname := res
	if len(gro) > 0 {
		schoolname = gro[0]
		log.Println("gro:", gro[0])
	}

	resp, err := itsl.CreateGroup(itswizard_basic.DbGroup15{
		SyncID:        p.SchuleUID,
		Name:          schoolname,
		ParentGroupID: "0",
	}, true)
	log.Println(schoolname)
	if err != nil {
		logContent = append(logContent, blusdLog{
			Success: true,
			Date:    time.Now(),
			Message: "School with id " + p.SchuleUID + " and name " + schoolname + " was not (!) created in itslearning. " + resp,
		})
	} else {
		logContent = append(logContent, blusdLog{
			Success: true,
			Date:    time.Now(),
			Message: "School with id " + p.SchuleUID + " and name " + schoolname + " was created in itslearning.",
		})
	}
	resp, err = itsl.CreateGroup(itswizard_basic.DbGroup15{
		SyncID:        p.SchuleUID + "-TEACHERS",
		Name:          "Lehrkräfte",
		ParentGroupID: p.SchuleUID,
	}, false)
	log.Println(p.SchuleUID + "-TEACHERS")

	if err != nil {
		logContent = append(logContent, blusdLog{
			Success: true,
			Date:    time.Now(),
			Message: "Group with id " + p.SchuleUID + "-TEACHERS " + "was not (!) created in itslearning. " + resp,
		})
	}
	resp, err = itsl.CreateGroup(itswizard_basic.DbGroup15{
		SyncID:        p.SchuleUID + "-STUDENTS",
		Name:          "Schüer",
		ParentGroupID: p.SchuleUID,
	}, false)
	log.Println(p.SchuleUID + "-STUDENTS")
	if err != nil {
		logContent = append(logContent, blusdLog{
			Success: true,
			Date:    time.Now(),
			Message: "Group with id " + p.SchuleUID + "-STUDENTS " + "was not (!) created in itslearning. " + resp,
		})
	}
	resp, err = itsl.CreateGroup(itswizard_basic.DbGroup15{
		SyncID:        p.SchuleUID + "-CLASSES",
		Name:          "Klassen",
		ParentGroupID: p.SchuleUID,
	}, false)
	log.Println(p.SchuleUID + "-CLASSES")
	if err != nil {
		logContent = append(logContent, blusdLog{
			Success: true,
			Date:    time.Now(),
			Message: "Group with id " + p.SchuleUID + "-CLASSES " + "was not (!) created in itslearning. " + resp,
		})
	}
	resp, err = itsl.CreateGroup(itswizard_basic.DbGroup15{
		SyncID:        p.SchuleUID + "-COURSES",
		Name:          "Kurse",
		ParentGroupID: p.SchuleUID,
	}, false)
	log.Println(p.SchuleUID + "-COURSES")
	if err != nil {
		logContent = append(logContent, blusdLog{
			Success: true,
			Date:    time.Now(),
			Message: "Group with id " + p.SchuleUID + "-COURSES " + "was not (!) created in itslearning. " + resp,
		})
	}
	return logContent, nil
}

func (p *LusdSchool) ExistInItslearning(itsl *imses.Request) (logContent []blusdLog, err error) {
	tmp := p.SchuleNummer + " " + p.SchuleName
	res := ""
	var gro []string
	for i, r := range tmp {
		res = res + string(r)
		if i > 0 && (i+1)%64 == 0 {
			gro = append(gro, res)
			res = ""
		}
	}
	schoolname := res
	if len(gro) > 0 {
		schoolname = gro[0]
		log.Println("gro:", gro[0])
	}

	resp, err := itsl.CreateGroup(itswizard_basic.DbGroup15{
		SyncID:        p.SchuleUID,
		Name:          schoolname,
		ParentGroupID: "0",
	}, true)
	log.Println(schoolname)
	if err != nil {
		logContent = append(logContent, blusdLog{
			Success: true,
			Date:    time.Now(),
			Message: "School with id " + p.SchuleUID + " and name " + schoolname + " was not (!) created in itslearning. " + resp,
		})
	} else {
		logContent = append(logContent, blusdLog{
			Success: true,
			Date:    time.Now(),
			Message: "School with id " + p.SchuleUID + " and name " + schoolname + " was created in itslearning.",
		})
	}
	resp, err = itsl.CreateGroup(itswizard_basic.DbGroup15{
		SyncID:        p.SchuleUID + "-TEACHERS",
		Name:          "Lehrkräfte",
		ParentGroupID: p.SchuleUID,
	}, false)
	log.Println(p.SchuleUID + "-TEACHERS")

	if err != nil {
		logContent = append(logContent, blusdLog{
			Success: true,
			Date:    time.Now(),
			Message: "Group with id " + p.SchuleUID + "-TEACHERS " + "was not (!) created in itslearning. " + resp,
		})
	}
	resp, err = itsl.CreateGroup(itswizard_basic.DbGroup15{
		SyncID:        p.SchuleUID + "-STUDENTS",
		Name:          "Schüer",
		ParentGroupID: p.SchuleUID,
	}, false)
	log.Println(p.SchuleUID + "-STUDENTS")
	if err != nil {
		logContent = append(logContent, blusdLog{
			Success: true,
			Date:    time.Now(),
			Message: "Group with id " + p.SchuleUID + "-STUDENTS " + "was not (!) created in itslearning. " + resp,
		})
	}
	resp, err = itsl.CreateGroup(itswizard_basic.DbGroup15{
		SyncID:        p.SchuleUID + "-CLASSES",
		Name:          "Klassen",
		ParentGroupID: p.SchuleUID,
	}, false)
	log.Println(p.SchuleUID + "-CLASSES")
	if err != nil {
		logContent = append(logContent, blusdLog{
			Success: true,
			Date:    time.Now(),
			Message: "Group with id " + p.SchuleUID + "-CLASSES " + "was not (!) created in itslearning. " + resp,
		})
	}
	resp, err = itsl.CreateGroup(itswizard_basic.DbGroup15{
		SyncID:        p.SchuleUID + "-COURSES",
		Name:          "Kurse",
		ParentGroupID: p.SchuleUID,
	}, false)
	log.Println(p.SchuleUID + "-COURSES")
	if err != nil {
		logContent = append(logContent, blusdLog{
			Success: true,
			Date:    time.Now(),
			Message: "Group with id " + p.SchuleUID + "-COURSES " + "was not (!) created in itslearning. " + resp,
		})
	}
	return logContent, nil
}
