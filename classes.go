package itswizard_m_berlinPreVersion

import (
	"encoding/json"
	itswizard_basic "github.com/itslearninggermany/itswizard_m_basic"
	imses "github.com/itslearninggermany/itswizard_m_imses"
	"log"
	"time"
)

type LusdClass struct {
	KlasseDtGeaendertAm string `json:"klasseDtGeaendertAm"`
	KlasseName          string `json:"klasseName"`
	KlasseUID           string `json:"klasseUID"`
	SchuleUID           string `json:"schuleUID"`
}

func classParse(out []byte, err error) (lusdClasses []LusdClass, err2 error) {
	if err != nil {
		return lusdClasses, err
	}
	err2 = json.Unmarshal(out, &lusdClasses)
	return lusdClasses, err2
}

/*
Returns all Schools from lusd
*/
func (p *BlusdConnection) GetAllClasses() (lusdClasses []LusdClass, out []byte, err error) {
	out, err = p.callAPI(time.Now(), "", class, false)
	if err != nil {
		return
	}
	lusdClasses, err = classParse(out, err)
	return
}

/*
Returns all new Courses created in the last 5 minutes
*/
func (p *BlusdConnection) GetAllClassessCreatedScinse(minutes time.Duration) (lusdClasses []LusdClass, err error) {
	newT := time.Now().Add(-time.Minute * minutes)
	log.Println(newT.Format(timeLayout))
	return classParse(p.callAPI(newT, "", class, true))
}

/*
Returns all Courses from a school
*/
func (p *BlusdConnection) GetAllClassesFromSchool(sid string) (lusdClasses []LusdClass, err error) {
	return classParse(p.callAPI(time.Now(), sid, class, false))
}

/*
Returns all new Courses from a school created in the last 5 minutes
*/
func (p *BlusdConnection) GetAllClassesFromSchoolScince(sid string, minutes time.Duration) (lusdClasses []LusdClass, err error) {
	newT := time.Now().Add(-time.Minute * minutes)
	log.Println(newT.Format(timeLayout))
	return classParse(p.callAPI(newT, sid, class, true))
}

/*
Returns all Courses from a school created since a given time
*/
func (p *BlusdConnection) GetAllClassesFromSchoolScinceAGivenTime(sid string, t time.Time) (lusdClasses []LusdClass, err error) {
	return classParse(p.callAPI(t, sid, class, true))
}

/*
Returns all Courses created since a given time
*/
func (p *BlusdConnection) GetAllClassesScinceAGivenTime(t time.Time) (lusdClasses []LusdClass, err error) {
	return classParse(p.callAPI(t, "", class, true))
}

/*
Imports all classes to itslearning. When the school exist it will be only updated.
*/
func (p *BlusdConnection) ImportAndUpdateAllClassesToItslearning(itsl *imses.Request) (logContent []blusdLog, err error) {
	lusdClasses, _, err := p.GetAllClasses()
	if err != nil {
		return logContent, err
	}
	for _, class := range lusdClasses {
		/*
			Todo: Schauen, ob es die Schule schon gibt.
		*/
		tmp := class.KlasseName
		res := ""
		var gro []string
		for i, r := range tmp {
			res = res + string(r)
			if i > 0 && (i+1)%64 == 0 {
				gro = append(gro, res)
				res = ""
			}
		}
		className := res
		if len(gro) > 0 {
			className = gro[0]
			log.Println("gro:", gro[0])
		}

		resp, err := itsl.CreateGroup(itswizard_basic.DbGroup15{
			SyncID:        class.KlasseUID,
			Name:          className,
			ParentGroupID: class.SchuleUID + "-CLASSES",
		}, false)
		log.Println(className)
		if err != nil {
			logContent = append(logContent, blusdLog{
				Success: true,
				Date:    time.Now(),
				Message: "Class with id " + class.KlasseUID + " and name " + className + " was not (!) created in itslearning. " + resp,
			})
		} else {
			logContent = append(logContent, blusdLog{
				Success: true,
				Date:    time.Now(),
				Message: "School with id " + class.KlasseUID + " and name " + className + " was created in itslearning.",
			})
		}
	}
	return logContent, nil
}

/*
TODO: Delete all classes which are not anymore in lusd.
*/
func (p *BlusdConnection) DeleteClassesInItslearning(itsl *imses.Request) (logContent []blusdLog, err error) {

	return logContent, nil
}
