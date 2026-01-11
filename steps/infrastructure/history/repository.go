package history

import (
    "context"
    "errors"
    "fmt"
    "os"
    "strconv"
    "time"

    "effective-architecture/steps/domain"
    "effective-architecture/steps/domain/history"
    "github.com/jackc/pgx/v5"
)

var (
    ErrCouldNotTruncate = errors.New("could not truncate label template")
    ErrCouldNotDelete   = errors.New("could not delete label template")
    ErrCouldNotCreate   = errors.New("could not create label template")
)

type Repository struct {
    conn *pgx.Conn
}

func NewRepository() (*Repository, error) {
    conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
    if err != nil {
        return nil, fmt.Errorf("unable to connect to database: %w", err)
    }

    return &Repository{
        conn: conn,
    }, nil
}

func (r Repository) Create(ctx context.Context, history history.History) error {
    sql := "INSERT INTO label_templates_history VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

    var (
        newManufacturerOrganizationName    *string
        newManufacturerOrganizationAddress *string
        newManufacturerEmail               *string
        newManufacturerSite                *string
    )

    if history.NewManufacturerOrganizationName != nil {
        newManufacturerOrganizationName = &history.NewManufacturerOrganizationName.Name
    }

    if history.NewManufacturerOrganizationAddress != nil {
        newManufacturerOrganizationAddress = &history.NewManufacturerOrganizationAddress.Address
    }

    if history.NewManufacturerEmail != nil {
        newManufacturerEmail = &history.NewManufacturerEmail.Value
    }

    if history.NewManufacturerSite != nil {
        newManufacturerSite = &history.NewManufacturerSite.Value
    }

    result, err := r.conn.Exec(ctx, sql, history.AggregateID.UUID.String(),
        strconv.Itoa(history.OrderKey), history.Action,
        newManufacturerOrganizationName, newManufacturerOrganizationAddress, newManufacturerEmail, newManufacturerSite,
        time.Now())
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return ErrCouldNotCreate
    }

    return nil
}

func (r Repository) List(ctx context.Context, aggregateID domain.LabelTemplateID) ([]history.History, error) {
    sql := `SELECT order_key, action, new_manufacturer_organization_name, new_manufacturer_organization_address,
       new_manufacturer_email, new_manufacturer_site FROM label_templates_history WHERE label_template_id=$1
                                                                                  ORDER BY order_key
   `

    rows, err := r.conn.Query(ctx, sql, aggregateID.UUID.String())
    if err != nil {
        return nil, fmt.Errorf("error select history list: %w", err)
    }

    result := make([]history.History, 0)

    for rows.Next() {
        var historyModel = History{}

        err := rows.Scan(&historyModel.OrderKey, &historyModel.Action,
            &historyModel.NewManufacturerOrganizationNameValue,
            &historyModel.NewManufacturerOrganizationAddressValue,
            &historyModel.NewManufacturerEmailValue, &historyModel.NewManufacturerSiteValue)
        if err != nil {
            return nil, fmt.Errorf("error scan history row: %w", err)
        }

        domainHistory, err := mapHistoryToDomain(historyModel, aggregateID)
        if err != nil {
            return nil, err
        }

        result = append(result, domainHistory)
    }

    return result, nil
}

func mapHistoryToDomain(historyModel History, aggregateID domain.LabelTemplateID) (history.History, error) {
    var (
        newManufacturerOrganizationName    *domain.OrganizationName
        newManufacturerOrganizationAddress *domain.OrganizationAddress
        newManufacturerEmail               *domain.Email
        newManufacturerSite                *domain.Site
    )

    if historyModel.NewManufacturerOrganizationNameValue != nil {
        newManufacturerOrganizationName2, err := domain.NewOrganizationName(
            *historyModel.NewManufacturerOrganizationNameValue)
        if err != nil {
            return history.History{}, err
        }

        newManufacturerOrganizationName = &newManufacturerOrganizationName2
    }

    if historyModel.NewManufacturerOrganizationAddressValue != nil {
        newManufacturerOrganizationAddress2, err := domain.NewOrganizationAddress(
            *historyModel.NewManufacturerOrganizationAddressValue)
        if err != nil {
            return history.History{}, err
        }

        newManufacturerOrganizationAddress = &newManufacturerOrganizationAddress2
    }

    if historyModel.NewManufacturerEmailValue != nil {
        newManufacturerEmail2, err := domain.NewEmail(*historyModel.NewManufacturerEmailValue)
        if err != nil {
            return history.History{}, err
        }

        newManufacturerEmail = &newManufacturerEmail2
    }

    if historyModel.NewManufacturerSiteValue != nil {
        newManufacturerSite2, err := domain.NewSite(*historyModel.NewManufacturerSiteValue)
        if err != nil {
            return history.History{}, err
        }

        newManufacturerSite = &newManufacturerSite2
    }

    domainHistory, err := history.NewHistory(aggregateID, historyModel.OrderKey, historyModel.Action,
        newManufacturerOrganizationName,
        newManufacturerOrganizationAddress, newManufacturerEmail,
        newManufacturerSite)
    if err != nil {
        return history.History{}, err
    }

    return domainHistory, nil
}

func (r Repository) Count(ctx context.Context, aggregateID domain.LabelTemplateID) (int, error) {
    sql := `SELECT count(*) FROM label_templates_history WHERE label_template_id=$1`

    var count int

    err := r.conn.QueryRow(ctx, sql, aggregateID.UUID.String()).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("failed to count history: %w", err)
    }

    return count, nil
}
